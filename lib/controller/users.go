/*
 * Copyright 2021 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"github.com/SENERGY-Platform/device-manager/lib/auth"
	permmodel "github.com/SENERGY-Platform/permission-search/lib/model"
)

func (this *Controller) deleteUserRights(token auth.Token, kind string, id string, userId string) (err error) {
	rights, err := this.com.GetPermissions(token, kind, id)
	if err != nil {
		return err
	}
	userrights := map[string]permmodel.Right{}
	for user, right := range rights.UserRights {
		if user != userId {
			userrights[user] = right
		}
	}
	rights.UserRights = userrights

	err = this.publisher.PublishRights(this.config.DeviceTopic, id, rights.ResourceRightsBase)
	if err != nil {
		return err
	}
	return nil
}

func (this *Controller) DeleteUser(userId string) error {
	token, err := auth.CreateToken("device-manager", userId)
	if err != nil {
		return err
	}
	admintoken, err := auth.CreateTokenWithRoles("device-manager", userId, []string{"admin", "system", "user"})
	if err != nil {
		return err
	}
	//devices
	devicesToDelete, userToDeleteFromDevices, err := this.com.ResourcesEffectedByUserDelete(token, this.config.DeviceTopic)
	if err != nil {
		return err
	}
	for _, id := range devicesToDelete {
		err = this.publisher.PublishDeviceDelete(id, userId)
		if err != nil {
			return err
		}
	}
	for _, id := range userToDeleteFromDevices {
		err = this.deleteUserRights(admintoken, this.config.DeviceTopic, id, userId)
		if err != nil {
			return err
		}
	}
	//device-groups
	deviceGroupToDelete, userToDeleteFromDeviceGroups, err := this.com.ResourcesEffectedByUserDelete(token, this.config.DeviceGroupTopic)
	if err != nil {
		return err
	}
	for _, id := range deviceGroupToDelete {
		err = this.publisher.PublishDeviceGroupDelete(id, userId)
		if err != nil {
			return err
		}
	}
	for _, id := range userToDeleteFromDeviceGroups {
		err = this.deleteUserRights(admintoken, this.config.DeviceGroupTopic, id, userId)
		if err != nil {
			return err
		}
	}
	//hubs
	hubToDelete, userToDeleteFromHubs, err := this.com.ResourcesEffectedByUserDelete(token, this.config.HubTopic)
	if err != nil {
		return err
	}
	for _, id := range hubToDelete {
		err = this.publisher.PublishHubDelete(id, userId)
		if err != nil {
			return err
		}
	}
	for _, id := range userToDeleteFromHubs {
		err = this.deleteUserRights(admintoken, this.config.HubTopic, id, userId)
		if err != nil {
			return err
		}
	}
	//locations
	locationToDelete, userToDeleteFromLocations, err := this.com.ResourcesEffectedByUserDelete(token, this.config.LocationTopic)
	if err != nil {
		return err
	}
	for _, id := range locationToDelete {
		err = this.publisher.PublishLocationDelete(id, userId)
		if err != nil {
			return err
		}
	}
	for _, id := range userToDeleteFromLocations {
		err = this.deleteUserRights(admintoken, this.config.LocationTopic, id, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
