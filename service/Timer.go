package service

import (
	log "github.com/wonderivan/logger"
	"sync"
	"ticket-service/database/model"
	"ticket-service/pkg/utils/redis"
	"time"
)

const redisKey = "orders_waiting_pay"
const expireDuration = 10 * time.Minute

// TimerFreeOrder 定时任务server
func (operator *ResourceOperator) TimerFreeOrder(timeAfter time.Duration, order model.UserOrder) {
	after := time.After(timeAfter)
	//compel := time.After(timeAfter + 2*time.Minute)
	field := order.UUID.String()
	go func() {
		for {
			select {
			case <-after:
				//执行流程: 从redis中读取若为空则return 若能拿到说明未支付，则删除订单并且释放库存。
				exists, err := redis.RedisClient.HExists(redisKey, field).Result()
				if err != nil {
					log.Error("定时任务.读取redis待支付订单失败 err:[%v]", err)
					return
				}
				if exists {
					tx, err := operator.TransactionBegin()
					if err != nil {
						log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
						return
					}
					defer func() {
						_ = tx.TransactionRollback()
					}()
					//订单过期后，释放库存，删除redis中键值。
					_, err = redis.RedisClient.Del(redisKey).Result()
					if err != nil {
						log.Error("定时任务.删除redis待支付订单失败 err:[%v]", err)
						return
					}
					var mutex sync.Mutex
					mutex.Lock()
					selector := make(map[string]interface{})
					selector[model.FieldScheduleID] = order.ScheduleID
					selector[model.FieldSeatType] = order.SeatType
					err = tx.Database.AddEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, model.FieldSeatNowNums, "1")
					if err != nil {
						log.Error("定时任务.恢复库存失败 err:[%v]", err)
						return
					}
					mutex.Unlock()
					err = tx.TransactionCommit()
					if err != nil {
						log.Error("定时任务.恢复库存 Commit 失败 err:[%v]", err)
						return
					}
					return
				} else {
					log.Error("定时任务.订单已取消：[%v] err:[%v]", field, err)
					return
				}

				//强制再删除一遍，防止删除失败 (没必要)
				//case <-compel:
				//	//执行流程: 从redis中读取若为空则return 若能拿到说明未支付，则删除订单并且释放库存。
				//	exists, err := redis.RedisClient.Exists(redisKey).Result()
				//	if err != nil {
				//		log.Error("超时强制释放任务.读取redis待支付订单失败 err:[%v]", err)
				//		return
				//	}
				//	if exists > 0 {
				//		tx, err := operator.TransactionBegin()
				//		if err != nil {
				//			log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
				//			return
				//		}
				//		defer func() {
				//			_ = tx.TransactionRollback()
				//		}()
				//		//订单过期后，释放库存，删除redis中键值。
				//		_, err = redis.RedisClient.HDel(redisKey, field).Result()
				//		if err != nil {
				//			log.Error("超时强制释放任务.删除redis待支付订单失败 err:[%v]", err)
				//			return
				//		}
				//		var mutex sync.Mutex
				//		mutex.Lock()
				//		selector := make(map[string]interface{})
				//		selector[model.FieldID] = order.ScheduleID
				//		selector[model.FieldSeatType] = order.SeatType
				//		err = tx.Database.AddEntityRowsByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, "seat_now_nums", "1")
				//		if err != nil {
				//			log.Error("超时强制释放任务.恢复库存失败 err:[%v]", err)
				//			return
				//		}
				//		mutex.Unlock()
				//		err = tx.TransactionCommit()
				//		if err != nil {
				//			log.Error("超时强制释放任务.恢复库存 Commit 失败 err:[%v]", err)
				//			return
				//		}
				//		return
				//	} else {
				//		log.Error("超时强制释放任务.订单非法删除！！！ err:[%v]", err)
				//		return
				//	}
			}
		}
	}()
}
