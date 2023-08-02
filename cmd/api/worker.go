package main

import (
	"context"
	"sync"
	"time"
)

func (app *application) verificationWorker(wg *sync.WaitGroup, ctx context.Context, waitTime int) {

	defer wg.Done()

	waitTimeDuration := time.Duration(waitTime) * time.Second
	ticker := time.NewTicker(waitTimeDuration)

	for {
		select {
		case <-ticker.C:
			select {
			case user := <-app.verificationQueue:
				user.VerificationStatus = true
			default:
				break
			}
		case <-ctx.Done():
			return
		default:
			continue
		}
	}

}

func (app *application) transactionWorker(wg *sync.WaitGroup, ctx context.Context, waitTime int) {

	defer wg.Done()

	waitTimeDuration := time.Duration(waitTime) * time.Second
	ticker := time.NewTicker(waitTimeDuration)

	for {
		select {
		case <-ticker.C:
			select {
			case trx := <-app.transactionQueue:
				sender, err := app.store.GetUserByID(trx.SenderID)
				if err != nil {
					break
				}
				receiver, err := app.store.GetUserByID(trx.ReceiverID)
				if err != nil {
					break
				}
				if sender.VerificationStatus == false {
					app.verificationQueue <- sender
				} else if receiver.VerificationStatus == false {
					app.verificationQueue <- receiver
				} else {
					app.transactionMutex.Lock()
					if trx.Amount <= sender.Balance {
						sender.Balance -= trx.Amount
						receiver.Balance += trx.Amount
					}
					app.transactionMutex.Unlock()
				}

			default:
				break
			}
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}
