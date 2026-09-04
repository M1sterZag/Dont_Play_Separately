package domain

import (
	"fmt"
	"hash/fnv"

	"github.com/google/uuid"
)

var DefaultAvatarCategories = []string{
	"action", "general", "moba", "neutral", "shooter", "sports_racing",
}

const DefaultAvatarsPerCategory = 9

func SelectDefaultAvatar(userID uuid.UUID) string {
	h := fnv.New64a()
	_, _ = h.Write(userID[:])
	idx := int(h.Sum64() % uint64(len(DefaultAvatarCategories)*DefaultAvatarsPerCategory))

	category := DefaultAvatarCategories[idx/DefaultAvatarsPerCategory]
	number := idx%DefaultAvatarsPerCategory + 1

	return fmt.Sprintf("%s/avatar_%s_%02d.png", category, category, number)
}

func IsValidAvatarKey(key string) bool {
	for _, category := range DefaultAvatarCategories {
		for n := 1; n <= DefaultAvatarsPerCategory; n++ {
			if key == fmt.Sprintf("%s/avatar_%s_%02d.png", category, category, n) {
				return true
			}
		}
	}

	return false
}
