package auth_service

import (
	"fmt"
	"hash/fnv"

	"github.com/google/uuid"
)

var defaultAvatarCategories = []string{
	"action", "general", "moba", "neutral", "shooter", "sports_racing",
}

const defaultAvatarsPerCategory = 9

func SelectDefaultAvatar(userID uuid.UUID) string {
	h := fnv.New64a()
	_, _ = h.Write(userID[:])
	idx := int(h.Sum64() % uint64(len(defaultAvatarCategories)*defaultAvatarsPerCategory))

	category := defaultAvatarCategories[idx/defaultAvatarsPerCategory]
	number := idx%defaultAvatarsPerCategory + 1

	return fmt.Sprintf("%s/avatar_%s_%02d.png", category, category, number)
}
