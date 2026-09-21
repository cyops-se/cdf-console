# Device Management

The Devices page provides a centralized view of all managed devices in your fleet.

## Device List

View all devices with the following information:

- **Status** - Online/offline indicator
- **Hostname** - Device identifier
- **Model** - Device type (RUTX, RUTM, etc.)
- **IP Address** - Current network address
- **Firmware** - Installed version
- **Last Seen** - Most recent contact time

## Device Actions

Available actions for each device:

- **Configure** - Modify device settings
- **Update Firmware** - See [Firmware Management](firmware-management.md)
- **View Logs** - Access syslog entries
- **SSH Access** - Open SSH session
- **Remove** - Delete from fleet

## Presence Monitoring

The system continuously monitors device availability using:

- ICMP ping checks (every 15 seconds)
- SSH session health validation (every 30 seconds)
- Traffic statistics collection

For monitoring details, see [Device Monitoring](monitoring.md).

## Related Topics

- [Firmware Management](firmware-management.md)
- [SSH Sessions](ssh-sessions.md)
- [Network Configuration](../networking/overview.md)
