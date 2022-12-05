package pkg

import (
	"context"
	"runtime"

	"github.com/puzpuzpuz/xsync"
	"github.com/remeh/sizedwaitgroup"
)

type userService struct {
	Wait  *sizedwaitgroup.SizedWaitGroup
	Count int
}

type UserService struct {
	*Server

	targets *xsync.MapOf[string, userService]
}

func NewUserService(s *Server) *UserService {

	return &UserService{
		Server:  s,
		targets: xsync.NewMapOf[userService](),
	}
}

// Note that WaitQueue is not using remote service. It is just native. Not confirmed with load balancing.
func (s *UserService) WaitQueue(u User, ctx context.Context) error {
	ut, ok := s.targets.Load(u.Token)
	if !ok {
		ut.Count = 0
		_ut := sizedwaitgroup.New(s.Config.QueueLimit)
		ut.Wait = &_ut
		runtime.KeepAlive(ut.Wait)
	}
	ut.Count++
	s.targets.Store(u.Token, ut)
	return ut.Wait.AddWithContext(ctx)

}

func (s *UserService) RemoveQueue(u User) {
	ut, ok := s.targets.Load(u.Token)
	if !ok {
		return
	}
	ut.Count--
	ut.Wait.Done()

	if ut.Count == 0 {
		s.targets.Delete(u.Token)
		ut.Wait = nil
		return
	}

	s.targets.Store(u.Token, ut)
}

// func (s *UserService) RemoveListener(u User) {
// 	t, ok := s.targets.Load(u.Token)
// 	if !ok {
// 		return
// 	}

// 	t.Count--
// 	if t.Count == 0 {
// 		s.targets.Delete(u.Token)
// 		close(t.Uchan) // ??
// 		return
// 	}
// 	s.targets.Store(u.Token, t)
// }

// func (s *UserService) listenUser(u User, ut chan UserRequest) {
// 	s.Logger.Debugf("user %s listening for queue", u.Token)
// 	listencontext, cancel := context.WithCancel(context.Background())
// 	defer s.Logger.Debugf("user %s listening done", u.Token)
// 	defer cancel()

// 	for v := range s.Database.ListenUserRequests(listencontext, u) {

// 		ut <- v

// 		if s.getQueue(u) == 0 {
// 			break
// 		}
// 	}

// }

// func (s *UserService) getQueue(u User) int {
// 	v, _ := s.targets.Load(u.Token)
// 	return v.Count
// }
