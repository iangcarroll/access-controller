# access-controller
Access Controller allows you to connect biometric readers from IDEMIA to Home Assistant, allowing modern access control without difficult wiring.

Normal access control systems require expensive controller boards, extensive wiring, and mechanical door locks, which are not feasible for residential installation. With Access Controller, the only wiring required is a single PoE+ cable to the IDEMIA reader mounted outside your home.

## Getting Started
You can run Access Controller on any device that has access to your Home Assistant server. Your IDEMIA device also needs to be able to talk to the Access Controller port.

```
# Install Access Controller
git clone https://github.com/iangcarroll/access-controller && cd access-controller

# Run the server
go run ./cmd/server -allowed-users 3000 -entity lock.abc_lock_mechanism -token abc.def.xyz -base-url http://homeassistant.local:8123
```

## Supported Devices
### Reader
Any IDEMIA reader that supports remote IP messages should work, including the SIGMA and MorphoWave lines. However, only the following has been tested to work:
* IDEMIA MorphoWave XP/SP/Compact on FW v2.13.3 in On-demand Security mode

### Locks
Any Home Assistant-paired lock that can be unlocked in Home Assistant can be controlled by Access Controller.

## Configuring your IDEMIA reader
You will need to enroll permitted users directly on the IDEMIA device and then pass them to `allowed-users`.

1. Upgrade your IDEMIA device to the latest firmware and consider resetting it to factory defaults.
2. Using MorphoBioToolBox (MBTB), go to Tools and select `Enable On-demand Security`. Enforced Security is not currently supported.
3. Once enabled, go to `Device Settings` -> `Controller Feedback` and enable `Feedback over IP`. Enable `Send Remote Message over IP` to the server and port that will host Access Controller. Disable `Host On No Response`.
4. Go to `Device Settings` -> `Events configuration` and ensure `Send to controller` is enabled for at least `User control successful`. You can also choose to send other events to Access Controller. If no events are enabled, we will not receive any data over the socket.
5. After users successfully authenticate, you should see messages flowing to your Access Controller.

## Security Considerations
Don't expose Access Controller to the internet, as it could receive spoofed messages from a rogue device. An attacker could also remove your IDEMIA reader externally and connect themselves to your network. Access Controller will only open your smart lock if the user ID is successfully recognized, so you should ensure your user IDs are not easily predicted. You can consider adding additional protections, such as detecting tamper events from the IDEMIA reader or the presence of rogue devices on your network.

Access Controller does not currently support Enforced Security. However, Enforced Security would allow you to validate the client sending messages.