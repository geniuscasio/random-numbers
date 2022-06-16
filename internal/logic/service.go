package logic

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"random-numbers/internal/adapters"
	"random-numbers/internal/common"
	"sort"
	"time"
)

type service struct {
	generator adapters.NumbersGenerator
	repo      adapters.UserPersistence
}

const generateStatsKey = "generate"

func NewService(generator adapters.NumbersGenerator, repo adapters.UserPersistence) Service {
	return &service{
		generator: generator,
		repo:      repo,
	}
}

func (s *service) IncrementStatistic(userID uuid.UUID) error {
	user, err := s.repo.ByID(userID)
	if err != nil {
		return err
	}

	user.CallsStatistic[generateStatsKey]++

	return s.repo.Update(*user)
}

func (s *service) LoginUser(email string, pwd string) (uuid.UUID, error) {
	hash := md5.Sum([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])

	userID, err := s.repo.ByCreds(email, hashedPwd)
	if err != nil {
		return uuid.Nil, err
	}

	return userID.ID, nil
}

func (s *service) CreateUser(user common.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.CallsStatistic = make(map[string]int64)

	hash := md5.Sum([]byte(user.Password))
	user.PasswordHash = hex.EncodeToString(hash[:])

	return s.repo.Create(user)
}

func (s *service) Generate(from, to, count int64, orderDesc bool) (*common.GenerateResponse, error) {
	numbers, err := s.generator.Get(from, to, count)
	if err != nil {
		return nil, err
	}

	return s.convertNumbersToResponse(numbers, orderDesc), nil
}

func (s *service) GetStatistic(userID uuid.UUID) (common.DetailsResponse, error) {
	user, err := s.repo.ByID(userID)
	if err != nil {
		return common.DetailsResponse{}, err
	}

	return common.DetailsResponse{
		Name:                user.Name,
		CreatedAt:           user.CreatedAt.String(),
		NumberOfGenerations: user.CallsStatistic[generateStatsKey],
	}, nil
}

func (s *service) convertNumbersToResponse(in []int64, orderDesc bool) *common.GenerateResponse {
	rows := make([]common.Row, 0, len(in))

	var counts = make(map[int64]int64, len(in))
	var order = make([]int64, 0, len(in))

	for _, number := range in {
		counts[number]++
	}

	for number := range counts {
		order = append(order, number)
	}

	s.sort(order, orderDesc)

	for _, number := range order {
		rows = append(rows, common.Row{
			Number: number,
			Count:  counts[number],
		})
	}

	res := common.GenerateResponse(rows)

	return &res
}

func (s *service) sort(in []int64, orderDesc bool) {
	sort.Slice(in, func(i, j int) bool {
		if orderDesc {
			return in[i] > in[j]
		}

		return in[j] > in[i]
	})
}
