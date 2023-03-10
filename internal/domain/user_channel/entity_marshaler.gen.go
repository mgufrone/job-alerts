// Code generated by mgufrone.dev/job-alerts/cmd/generate-domain, DO NOT EDIT.
package user_channel

import (
	"encoding/json"
	jsonparser "github.com/buger/jsonparser"
	uuid "github.com/google/uuid"
	try "github.com/mgufrone/go-utils/try"
	user "mgufrone.dev/job-alerts/internal/domain/user"
	helpers "mgufrone.dev/job-alerts/pkg/helpers"
)

func (e *Entity) UnmarshalJSON(data []byte) error {
	return try.RunOrError(func() error {
		val, err := helpers.GetUUID(data, "id")
		if err != nil {
			return err
		}
		return e.SetID(val)
	}, func() error {
		var ref user.Entity
		val, err := helpers.GetUUID(data, "userID")
		if err != nil {
			return err
		}
		if val == uuid.Nil {
			return nil
		}
		roles, err := jsonparser.GetInt(data, "user_roles")
		if err != nil {
			return err
		}
		_ = ref.SetRoles(user.Role(roles))
		_ = ref.SetID(val)
		return e.SetUser(&ref)
	}, func() error {
		val, err := jsonparser.GetString(data, "channelType")
		if err != nil {
			return err
		}
		return e.SetChannelType(val)
	}, func() error {
		val, err := jsonparser.GetString(data, "receiver")
		if err != nil {
			return err
		}
		return e.SetReceiver(val)
	}, func() error {
		val, err := jsonparser.GetBoolean(data, "isActive")
		if err != nil {
			return err
		}
		return e.SetIsActive(val)
	}, func() error {
		val, err := helpers.GetTime(data, "createdAt")
		if err != nil {
			return err
		}
		return e.SetCreatedAt(val)
	}, func() error {
		val, err := helpers.GetTime(data, "updatedAt")
		if err != nil {
			return err
		}
		return e.SetUpdatedAt(val)
	})
}
func (e *Entity) MarshalJSON() ([]byte, error) {
	res := map[string]interface{}{
		"id": e.ID(),
		"user_roles": func() int {
			return int(e.User().Roles())
		}(),
		"userID": func() uuid.UUID {
			if e.User() == nil {
				return uuid.Nil
			}
			return e.User().ID()
		}(),
		"channelType": e.ChannelType(),
		"receiver":    e.Receiver(),
		"isActive":    e.IsActive(),
		"createdAt":   e.CreatedAt().UnixMilli(),
		"updatedAt":   e.UpdatedAt().UnixMilli(),
	}
	return json.Marshal(res)
}
