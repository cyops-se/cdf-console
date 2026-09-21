# Device Dashboard
>**Version**: 1.1<br>
>**Build**: 342

The Device Dashboard provides a real-time overview of all managed devices in your fleet.

![Sample Dashboard Icon](images/sample-icon.png)

> **Note:** The image above is embedded in the Go binary and served via `/api/help/images/sample-icon.png`. You can add your own screenshots and diagrams to the `server/help/docs/images/` folder.

## Dashboard Overview

The dashboard displays device cards with the following information:

### Device Status
- **Green Card** - Device is online and reachable
- **Red Card** - Device is offline or unreachable
- **Status Chip** - Shows "reachable" or "unreachable" status

### Device Information
- **Device Name/Hostname** - Primary identifier for the device
- **Serial Number** - Device serial number
- **Last Check** - Timestamp of the most recent status check

## Interface Statistics

Each device card shows network interface statistics when available:

### Interface Types
- **WAN** - Wide Area Network interface (internet connection)
- **LAN** - Local Area Network interface (internal network)
- **APN** - Mobile/4G cellular interface

### Statistics Displayed

For each active interface:

- **Operational State** - Interface status (up/down)
- **RX (Receive)** - Data received on the interface
- **TX (Transmit)** - Data transmitted on the interface
- **Packet Counts** - Number of packets received/transmitted

### Data Format

Traffic is displayed in human-readable format:
- B (Bytes)
- KB (Kilobytes)
- MB (Megabytes)
- GB (Gigabytes)
- TB (Terabytes)

Packet counts are abbreviated:
- K = Thousands (1000)
- M = Millions (1,000,000)

## Dashboard Actions

### Refresh
Click the refresh icon to manually update device statistics. The dashboard automatically refreshes every 30 seconds.

### Reset Statistics
Click the counter icon to reset all traffic statistics for all devices. This will:
- Clear all RX/TX byte counters
- Reset packet counters
- Start fresh statistics collection

**Note:** This action cannot be undone. Use with caution.

## Auto-Refresh

The dashboard automatically refreshes every 30 seconds to provide near real-time monitoring of your device fleet. You'll see:
- Updated interface statistics
- Current device status
- Latest last check timestamps

## Troubleshooting

### No Interface Statistics

If a device card shows "No interface statistics available":
- Device may be unreachable
- Network interfaces may not be configured
- Statistics collection may not have started yet

### Device Shows as Unreachable

If a device appears offline:
1. Check network connectivity to the device
2. Verify device is powered on
3. Check device IP address configuration
4. See [Troubleshooting](../troubleshooting/common-issues.md) for more help

## Related Topics

- [Device Management](../devices/overview.md)
- [Device Monitoring](../devices/monitoring.md)
- [Network Configuration](../networking/overview.md)
- [Troubleshooting](../troubleshooting/common-issues.md)
