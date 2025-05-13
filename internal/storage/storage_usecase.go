package storage

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/common"
	"arctfrex-customers/internal/common/enums"
	"arctfrex-customers/internal/model"
	"arctfrex-customers/internal/repository"
	user_mobile "arctfrex-customers/internal/user/mobile"
)

type StorageUsecase interface {
	UploadFile(userId, accountId, documentType, contentType string, file multipart.File, fileName string, fileSize int64) error
	DownloadFile(fileName string) ([]byte, error)
	PreviewFile(fileName string) ([]byte, error)
}

type storageUsecase struct {
	storageMinioClient MinioClient
	userRepository     user_mobile.UserRepository
	depositRepository  repository.DepositRepository
	accountRepository  repository.AccountRepository
}

func NewStorageUsecase(
	mc MinioClient,
	ur user_mobile.UserRepository,
	dr repository.DepositRepository,
	ar repository.AccountRepository,
) *storageUsecase {
	return &storageUsecase{
		storageMinioClient: mc,
		userRepository:     ur,
		depositRepository:  dr,
		accountRepository:  ar,
	}
}

func (su *storageUsecase) UploadFile(userId, accountId, documentType, contentType string, file multipart.File, fileName string, fileSize int64) error {
	if strings.ToLower(documentType) == "deposit" {
		depositPending, _ := su.depositRepository.GetPendingAccountByAccountIdUserId(accountId, userId)
		if depositPending != nil && depositPending.IsActive {
			log.Println("deposit still in pending approval")
			return errors.New("deposit still in pending approval")
		}
	}

	pathUrl := "https://" + os.Getenv(common.MINIO_BASEURL) + "/" + os.Getenv(common.MINIO_BUCKET_NAME)
	fileName = time.Now().Format("20060102150405") + "_" + documentType + "_" + fileName
	filePath := pathUrl + "/" + fileName
	fmt.Println(contentType)
	err := su.storageMinioClient.UploadFile(contentType, file, fileName, fileSize)
	if err != nil {
		log.Println(err)
		return err
	}

	switch strings.ToLower(documentType) {
	case "ktp":
		{
			return su.userRepository.UpdateProfileKtpPhoto(&user_mobile.UserProfile{ID: userId, KtpPhoto: filePath})
		}
	case "selfie":
		{
			return su.userRepository.UpdateProfileSelfiePhoto(&user_mobile.UserProfile{ID: userId, SelfiePhoto: filePath})
		}
	case "npwp":
		{
			return su.userRepository.UpdateProfileNpwpPhoto(&user_mobile.UserProfile{ID: userId, NpwpPhoto: filePath})
		}
	case "declaration":
		{
			return su.userRepository.UpdateProfileDeclarationVideo(&user_mobile.UserProfile{ID: userId, DeclarationVideo: filePath})
		}
	case "deposit":
		{
			depositdb, _ := su.depositRepository.GetNewDepositByAccountIdUserId(accountId, userId)

			if depositdb == nil {
				depositID, err := uuid.NewUUID()
				if err != nil {
					return err
				}

				depositdb = &model.Deposit{
					ID:             common.UUIDNormalizer(depositID),
					AccountID:      accountId,
					UserID:         userId,
					ApprovalStatus: enums.DepositApprovalStatusNew,
					BaseModel:      base.BaseModel{IsActive: true},
				}
			}

			depositdb.DepositPhoto = filePath

			return su.depositRepository.SaveDepositPhoto(depositdb)
		}
	case "realaccount_callrecording":
		{
			return su.accountRepository.UpdateRealAccountCallRecording(&model.Account{ID: accountId, UserID: userId, RealAccountCallRecording: filePath})
		}
	default:
		{
			return su.userRepository.UpdateProfileAdditionalDocumentPhoto(&user_mobile.UserProfile{ID: userId, AdditionalDocumentPhoto: filePath})
		}
	}

	// if strings.ToLower(documentType) == "ktp" {
	// 	return su.userRepository.UpdateProfileKtpPhoto(&user_mobile.UserProfile{ID: userId, KtpPhoto: filePath})
	// }

	// if strings.ToLower(documentType) == "selfie" {
	// 	return su.userRepository.UpdateProfileSelfiePhoto(&user_mobile.UserProfile{ID: userId, SelfiePhoto: filePath})
	// }

	// if strings.ToLower(documentType) == "npwp" {
	// 	return su.userRepository.UpdateProfileNpwpPhoto(&user_mobile.UserProfile{ID: userId, NpwpPhoto: filePath})
	// }

	// if strings.ToLower(documentType) == "additional_document" {
	// 	return su.userRepository.UpdateProfileAdditionalDocumentPhoto(&user_mobile.UserProfile{ID: userId, AdditionalDocumentPhoto: filePath})
	// }

	// if strings.ToLower(documentType) == "declaration" {
	// 	return su.userRepository.UpdateProfileDeclarationVideo(&user_mobile.UserProfile{ID: userId, DeclarationVideo: filePath})
	// }

	// if strings.ToLower(documentType) == "deposit" {
	// 	depositdb, _ := su.depositRepository.GetNewDepositByAccountIdUserId(accountId, userId)

	// 	if depositdb == nil {
	// 		depositID, err := uuid.NewUUID()
	// 		if err != nil {
	// 			return err
	// 		}

	// 		depositdb = &deposit.Deposit{
	// 			ID:        common.UUIDNormalizer(depositID),
	// 			AccountID: accountId,
	// 			UserID:    userId,
	// 			// DepositPhoto: filePath,
	// 			BaseModel: base.BaseModel{IsActive: true},
	// 		}
	// 	}

	// 	depositdb.DepositPhoto = filePath

	// 	return su.depositRepository.SaveDepositPhoto(depositdb)
	// }

	//return nil
}

func (su *storageUsecase) DownloadFile(fileName string) ([]byte, error) {
	return su.storageMinioClient.DownloadFile(fileName)
}

func (su *storageUsecase) PreviewFile(fileName string) ([]byte, error) {
	return su.storageMinioClient.DownloadFile(fileName)
}
