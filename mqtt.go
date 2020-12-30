/*
 * talkkonnect headless mumble client/gateway with lcd screen and channel control
 * Copyright (C) 2018-2019, Suvir Kumar <suvir@talkkonnect.com>
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * talkkonnect is the based on talkiepi and barnard by Daniel Chote and Tim Cooper
 *
 * The Initial Developer of the Original Code is
 * Suvir Kumar <suvir@talkkonnect.com>
 * Portions created by the Initial Developer are Copyright (C) Suvir Kumar. All Rights Reserved.
 *
 * talKKonnectContributor(s):
 *
 * Suvir Kumar <suvir@talkkonnect.com>
 *
 * My Blog is at www.talkkonnect.com
 * The source code is hosted at github.com/talkkonnect
 *
 * MQTT License Details Copyright (c) 2013 IBM Corp.
 *
 * This project is dual licensed under the Eclipse Public License 1.0 and the
 * Eclipse Distribution License 1.0 as described in the epl-v10 and edl-v10 files.
 * The EDL is copied below in order to pass the pkg.go.dev license check (https://pkg.go.dev/license-policy).
 * Eclipse Distribution License - v 1.0
 * Copyright (c) 2007, Eclipse Foundation, Inc. and its licensors.
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 * Neither the name of the Eclipse Foundation, Inc. nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 * All rights reserved. This program and the accompanying materials are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at http://www.eclipse.org/legal/epl-v10.html
 *
 */

package talkkonnect

import (
	"crypto/tls"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/talkkonnect/gpio"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type dateTimeScheduleStruct struct {
	startDateTime string
	endDateTime   string
	matched       bool
	defaultLogic  bool
	stopOnMatch   bool
}

type dayScheduleStruct struct {
	dayint       int
	startTime    int
	endTime      int
	matched      bool
	defaultLogic bool
	stopOnMatch  bool
}

func relayAllPulse() {
	if relayAllState {
		relayCommand(0, "off")
	} else {
		relayCommand(0, "on")
	}
	relayAllState = !relayAllState
}

func relayCommand(relayNo int, command string) {
	// all relays (0)
	if relayNo == 0 {
		for i := 1; i <= int(TotalRelays); i++ {
			if command == "on" {
				log.Println("info: Relay ", i, "On")
				gpio.NewOutput(RelayPins[i], false)

			}
			if command == "off" {
				log.Println("info: Relay ", i, "Off")
				gpio.NewOutput(RelayPins[i], true)
			}
			if command == "pulse" {
				log.Println("info: Relay ", i, "Pulse")
				gpio.NewOutput(RelayPins[i], false)
				time.Sleep(RelayPulseMills * time.Millisecond)
				gpio.NewOutput(RelayPins[i], true)
			}
		}
		return
	}

	//specific relay (Number Between 1 and TotalRelays)
	if relayNo >= 0 && relayNo <= int(TotalRelays) {
		if command == "on" {
			log.Println("info: Relay ", relayNo, "On")
			gpio.NewOutput(RelayPins[relayNo], false)
		}
		if command == "off" {
			log.Println("info: Relay ", relayNo, "Off")
			gpio.NewOutput(RelayPins[relayNo], true)
		}
		if command == "pulse" {
			log.Println("info: Relay ", relayNo, "Pulse")
			gpio.NewOutput(RelayPins[relayNo], false)
			time.Sleep(RelayPulseMills * time.Millisecond)
			gpio.NewOutput(RelayPins[relayNo], true)
		}
	}
}

func dateTimeWithinRange(dateTimeSchedule dateTimeScheduleStruct) (bool, bool, bool, error) {
	var dateFormat string = "02/01/2006 15:04"
	startDateTime, err := time.Parse(dateFormat, dateTimeSchedule.startDateTime)
	if err != nil {
		return false, false, false, err
	}

	endDateTime, err := time.Parse(dateFormat, dateTimeSchedule.endDateTime)
	if err != nil {
		return false, false, false, err
	}

	checkDateTime, err := time.Parse(dateFormat, time.Now().Format("02/01/2006 15:04"))
	if err != nil {
		return false, false, false, err
	}
	log.Println("------")
	log.Println("debug: startdate is ", startDateTime, " enddate is ", endDateTime, " check date is ", checkDateTime)
	if startDateTime.Before(checkDateTime) && endDateTime.After(checkDateTime) {
		return true, dateTimeSchedule.defaultLogic, dateTimeSchedule.stopOnMatch, nil
	}
	return false, dateTimeSchedule.defaultLogic, dateTimeSchedule.stopOnMatch, nil
}

//func dayTimeWithinRange(startTime string, endTime string, dayCheck string, dateFormat string, defaultLogicDay string) (bool, error) {
func dayTimeWithinRange(dayTimeWithinRange dayScheduleStruct) (bool, bool, bool, error) {

	t1 := time.Now()
	t1Day := int(t1.Weekday())
	t1Minute := int((t1.Hour() * 60) + t1.Minute())

	log.Println("------")
	log.Println("debug: day is ", t1Day, " starttime is ", dayTimeWithinRange.startTime, " endtime is ", dayTimeWithinRange.endTime, " checkday is ", t1Day, " check time is ", t1Minute)

	if t1Day == dayTimeWithinRange.dayint && (dayTimeWithinRange.startTime <= t1Minute && dayTimeWithinRange.endTime >= t1Minute) {
		return true, dayTimeWithinRange.defaultLogic, dayTimeWithinRange.stopOnMatch, nil
	}
	return false, dayTimeWithinRange.defaultLogic, dayTimeWithinRange.stopOnMatch, nil
}

func mqtttestpub() {

	if MQTTAction != "pub" {
		log.Println("error: Invalid setting for -action, must be pub for test pub")
		return
	}

	if MQTTTopic == "" {
		log.Println("error: Invalid setting for -topic, must not be empty")
		return
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(MQTTBroker)
	opts.SetClientID(MQTTId)
	opts.SetUsername(MQTTUser)
	opts.SetPassword(MQTTPassword)
	opts.SetCleanSession(MQTTCleansess)

	if MQTTStore != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(MQTTStore))
	}

	if MQTTAction == "pub" {

		log.Printf("info: action      : %s\n", MQTTAction)
		log.Printf("info: broker      : %s\n", MQTTBroker)
		log.Printf("info: clientid    : %s\n", MQTTId)
		log.Printf("info: user        : %s\n", MQTTUser)
		log.Printf("info: mqttpassword: %s\n", MQTTPassword)
		log.Printf("info: topic       : %s\n", MQTTTopic)
		log.Printf("info: message     : %s\n", MQTTPayload)
		log.Printf("info: qos         : %d\n", MQTTQos)
		log.Printf("info: cleansess   : %v\n", MQTTCleansess)
		log.Printf("info: num         : %d\n", MQTTNum)
		log.Printf("info: store       : %s\n", MQTTStore)

		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		log.Println("info: Test MQTT Publisher Started")
		for i := 0; i < MQTTNum; i++ {
			log.Println("info: Publishing MQTT Message")
			token := client.Publish(MQTTTopic, byte(MQTTQos), false, MQTTPayload)
			token.Wait()
		}

		client.Disconnect(250)

		log.Println("info: Test MQTT Publisher Disconnected")
	}
}

func (b *Talkkonnect) mqttsubscribe() {

	log.Printf("info: MQTT Subscription Information")
	log.Printf("info: MQTT Broker      : %s\n", MQTTBroker)
	log.Printf("debug: MQTT clientid    : %s\n", MQTTId)
	log.Printf("debug: MQTT user        : %s\n", MQTTUser)
	log.Printf("debug: MQTT password    : %s\n", MQTTPassword)
	log.Printf("info: Subscribed topic : %s\n", MQTTTopic)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	connOpts := MQTT.NewClientOptions().AddBroker(MQTTBroker).SetClientID(MQTTId).SetCleanSession(true)
	if MQTTUser != "" {
		connOpts.SetUsername(MQTTUser)
		if MQTTPassword != "" {
			connOpts.SetPassword(MQTTPassword)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(MQTTTopic, byte(MQTTQos), b.onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		log.Printf("info: Connected to     : %s\n", MQTTBroker)
	}

	<-c
}

func (b *Talkkonnect) onMessageReceived(client MQTT.Client, message MQTT.Message) {
	log.Printf("info: Received MQTT message on topic: %s Payload: %s\n", message.Topic(), message.Payload())

	switch string(message.Payload()) {
	case "DEL":
		log.Println("info: MQTT Display Menu Request Processed Succesfully")
		b.commandKeyDel()
	case "F1":
		log.Println("API Channel Up Request Processed Succesfully\n")
		b.commandKeyF1()
	case "F2":
		log.Println("API Channel Down Request Processed Succesfully\n")
		b.commandKeyF2()
	case "F3":
		log.Println("API Mute/UnMute Speaker Request Processed Succesfully\n")
		b.commandKeyF3("toggle")
	case "F3-mute":
		log.Println("API Mute/UnMute Speaker Request Processed Succesfully\n")
		b.commandKeyF3("mute")
	case "F3-unmute":
		log.Println("API Mute/UnMute Speaker Request Processed Succesfully\n")
		b.commandKeyF3("unmute")
	case "F4":
		log.Println("API Current Volume Level Request Processed Succesfully\n")
		b.commandKeyF4()
	case "F5":
		log.Println("API Digital Volume Up Request Processed Succesfully\n")
		b.commandKeyF5()
	case "F6":
		log.Println("API Digital Volume Down Request Processed Succesfully\n")
		b.commandKeyF6()
	case "F7":
		log.Println("API List Server Channels Request Processed Succesfully\n")
		b.commandKeyF7()
	case "F8":
		log.Println("API Start Transmitting Request Processed Succesfully\n")
		b.commandKeyF8()
	case "F9":
		log.Println("API Stop Transmitting Request Processed Succesfully\n")
		b.commandKeyF9()
	case "F10":
		log.Println("API List Online Users Request Processed Succesfully\n")
		b.commandKeyF10()
	case "F11":
		log.Println("API Play/Stop Chimes Request Processed Succesfully\n")
		b.commandKeyF11()
	case "F12":
		log.Println("API Request GPS Position Processed Succesfully\n")
		b.commandKeyF12()
	case "commandKeyCtrlE":
		log.Println("API Send Email Processed Succesfully\n")
		b.commandKeyCtrlE()
	case "commandKeyCtrlF":
		log.Println("API Previous Server Processed Successfully\n")
		b.commandKeyCtrlF()
	case "commandKeyCtrlN":
		log.Println("API Next Server Processed Successfully\n")
		b.commandKeyCtrlN()
	case "commandKeyCtrlL":
		log.Println("API Clear Screen Processed Successfully\n")
		b.commandKeyCtrlL()
	case "commandKeyCtrlO":
		log.Println("API Ping Servers Processed Succesfully\n")
		b.commandKeyCtrlO()
	case "commandKeyCtrlP":
		log.Println("API Request Panic Simulation Processed Succesfully\n")
		b.commandKeyCtrlP()
	case "commandKeyCtrlR":
		log.Println("API Request Repeat Tx Loop Test Processed Succesfully\n")
		b.commandKeyCtrlR()
	case "commandKeyCtrlS":
		log.Println("API Request Scan Processed Succesfully\n")
		b.commandKeyCtrlS()
	case "commandKeyCtrlT":
		log.Println("API Request Show Acknowledgements Processed Succesfully\n")
		b.commandKeyCtrlT()
	case "commandKeyCtrlV":
		log.Println("API Request Current Version Succesfully\n")
		b.commandKeyCtrlV()
	case "commandKeyCtrlX":
		log.Println("API Print XML Config Processed Succesfully\n")
		b.commandKeyCtrlX()
	case "attentionledon":
		log.Println("API Turn On Attention LED Succesfully\n")
  		b.LEDOn(b.AttentionLED)
	case "attentionledoff":
		log.Println("API Turn Off Attention LED Succesfully\n")
  		b.LEDOff(b.AttentionLED)
	// todo add other automation control for buttons, relays and leds here as needed in the future
	default:
		log.Printf("error: Undefined Command Received MQTT message on topic: %s Payload: %s\n", message.Topic(), message.Payload())
	}
	return
}
