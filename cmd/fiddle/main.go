package main

//func main() {
//	ctx := context.Background()
//	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
//	defer cancel()
//
//	stream := &fsevents.EventStream{
//		Paths:   []string{"/"},
//		Flags:   fsevents.FileEvents,
//		EventID: 0,
//		//Resume:  true,
//	}
//	stream.Start()
//
//	done := make(chan struct{})
//	go func() {
//		<-ctx.Done()
//		stream.Stop()
//		stream.Flush( /* sync */ true)
//		close(done)
//	}()
//
//loop:
//	for {
//		select {
//		case events := <-stream.Events:
//			for _, event := range events {
//				fmt.Printf("%#v\n", event)
//			}
//		case <-done:
//			break loop
//		}
//	}
//}
