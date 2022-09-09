package endpoints

import (
	"math/rand"
	"sync"
	"time"
)

type authStateManager struct {
	rand        *rand.Rand
	knownStates map[string]*authState
	lock        sync.RWMutex
}

type authState struct {
	ExpiresAt time.Time
	NextURL   string
	Key       string
}

func newAuthStateManager() *authStateManager {
	asm := &authStateManager{
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
		knownStates: make(map[string]*authState),
	}
	go asm.authStateGarbageCollector()
	return asm
}

func (asm *authStateManager) authStateGarbageCollector() {
	for {
		time.Sleep(time.Minute * 5)
		asm.lock.Lock()

		var toDelete []string

		for k, v := range asm.knownStates {
			if !v.ExpiresAt.After(time.Now().UTC()) {
				toDelete = append(toDelete, k)
			}
		}

		for _, k := range toDelete {
			delete(asm.knownStates, k)
		}

		asm.lock.Unlock()
	}
}

func (asm *authStateManager) Get(key string) *authState {
	asm.lock.RLock()
	defer asm.lock.RUnlock()

	if v, found := asm.knownStates[key]; found && v.ExpiresAt.After(time.Now().UTC()) {
		return v
	}
	return nil
}

func (asm *authStateManager) Delete(key string) {
	asm.lock.Lock()
	defer asm.lock.Unlock()

	delete(asm.knownStates, key)
}

func (asm *authStateManager) New(nextURL string) *authState {
	const timeout = time.Minute * 2

	b := make([]byte, 30)
	for i := 0; i < len(b); i++ {
		b[i] = byte('A' + rand.Intn(25))
	}

	asm.lock.Lock()
	defer asm.lock.Unlock()

	// Yes - there is a very small chance of collision here. However, if it
	// does happen, all it'll do is prevent a single user's login from
	// working once. They can just retry and almost certainly be fine, so we'll
	// just ignore that whole problem.

	as := &authState{
		ExpiresAt: time.Now().UTC().Add(timeout),
		NextURL:   nextURL,
		Key:       string(b),
	}

	asm.knownStates[string(b)] = as

	return as
}
