/*
 * @Author: Jiaqin,Lu
 * @Date: 2025-02-20 02:47:28
 * @Last Modified by: Jiaqin,Lu
 * @Last Modified time: 2025-02-20 03:09:53
 */
package service

import (
	"context"
	"strconv"
	"time"

	"github.com/PiaoAdmin/gomall/app/payment/biz/dal/mysql"
	"github.com/PiaoAdmin/gomall/app/payment/biz/model"
	payment "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/kitex/pkg/kerrors"
	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
)

type ChargeService struct {
	ctx context.Context
} // NewChargeService new ChargeService
func NewChargeService(ctx context.Context) *ChargeService {
	return &ChargeService{ctx: ctx}
}

// Run create note info
func (s *ChargeService) Run(req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	// Finish your business logic.
	//暂时只提供银行卡号一种方式
	card := creditcard.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CreditCard.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CreditCard.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CreditCard.CreditCardExpirationYear)),
	}
	//验证银行卡有效性
	err = card.Validate(true)
	if err != nil {
		return nil, kerrors.NewBizStatusError(400, err.Error())
	}

	//使用随机生成transcationId，在之后调用真实API获得
	transcationId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	//将信息插入数据库
	err = model.CreatePaymentLog(mysql.DB, s.ctx, &model.PaymentLog{
		UserId: int64(req.UserId),
		OrderId: func() int64 {
			id, _ := strconv.ParseInt(req.OrderId, 10, 64)
			return id
		}(),
		TranscationId: func() int64 {
			id, _ := strconv.ParseInt(transcationId.String(), 10, 64)
			return id
		}(),
		Amount: req.Amount,
		PayAt:  time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &payment.ChargeResp{TransactionId: transcationId.String()}, nil
}
