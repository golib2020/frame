package queue

//TestQueuePop 暂时保留记录
//func TestQueuePop(t *testing.T) {
//
//	topic := "default"
//	redisQueue := NewRedisQueue(f.Redis(), time.Minute*10)
//	for {
//		//time.Sleep(time.Millisecond * 100)
//		job := NewJob("")
//		if err := redisQueue.Pop(topic, job); err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		data := job.GetBody()
//
//		if data == "" {
//			continue
//		}
//
//		go func() {
//			defer func() {
//				if p := recover(); p != nil { //执行失败进行重试
//					log.Printf("panic recover! p: %v\n", p)
//					debug.PrintStack()
//					redisQueue.Release(topic, job, time.Minute)
//				} else { //执行成功后直接删除任务
//					if err := redisQueue.Delete(topic, job); err != nil {
//						fmt.Println(err)
//					}
//				}
//			}()
//
//			//正常业务逻辑
//			fmt.Println(data)
//
//		}()
//	}
//}
//
//func TestQueuePush(t *testing.T) {
//	topic := "default"
//
//	redisQueue := NewRedisQueue(f.Redis(), time.Minute*10)
//
//	redisQueue.Push(topic, NewJob("0s"))
//
//	//for i := 0; i < 10000; i++ {
//	//	redisQueue.Push(topic, NewJob(fmt.Sprintf("正常推送%5d", i)))
//	//}
//
//	redisQueue.Later(topic, NewJob("延时推送10s"), time.Second*10)
//	redisQueue.Later(topic, NewJob("延时推送1分钟"), time.Minute)
//	redisQueue.Later(topic, NewJob("延时推送1小时"), time.Hour)
//
//	fmt.Println("完成")
//
//}
