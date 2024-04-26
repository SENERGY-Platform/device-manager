/*
 * Copyright 2019 InfAI (CC SES)
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

package model

type DeviceUpdateOptions struct {
	UpdateOnlySameOriginAttributes []string
	Wait                           bool
}

type DeviceCreateOptions struct {
	Wait bool
}

type DeviceDeleteOptions struct {
	Wait bool
}

type DeviceTypeUpdateOptions struct {
	Wait bool
}

type DeviceTypeDeleteOptions struct {
	Wait bool
}

type HubUpdateOptions struct {
	Wait bool
}

type HubDeleteOptions struct {
	Wait bool
}

type AspectUpdateOptions struct {
	Wait bool
}

type AspectDeleteOptions struct {
	Wait bool
}
