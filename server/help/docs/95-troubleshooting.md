# Felsökning

This guide covers frequently encountered issues and their solutions.

## Device Not Responding

**Symptoms:** Device shows as offline, ping fails

**Solutions:**
1. Verify network connectivity
2. Check device IP address configuration
3. Ensure device is powered on
4. Review firewall rules blocking ICMP

See [Network Issues](network-issues.md) for detailed troubleshooting.

## SSH Connection Failed

**Symptoms:** Cannot establish SSH session

**Solutions:**
1. Verify SSH credentials
2. Check if SSH service is running on device
3. Review SSH session status at `/api/device/ssh/sessions`
4. Try closing and reopening session

See [SSH Sessions](../devices/ssh-sessions.md) for session management.

## Firmware Update Failed

**Symptoms:** Firmware upgrade process errors

**Solutions:**
1. Verify firmware file integrity
2. Ensure sufficient storage space on device
3. Check device model matches firmware
4. Review syslog for error details

See [Firmware Management](../devices/firmware-management.md) for update procedures.

## Onboarding Stuck

**Symptoms:** Onboarding process hangs at specific step

**Solutions:**
1. Check device accessibility
2. Verify network configuration
3. Review onboarding logs
4. Restart onboarding process

See [Device Onboarding](../getting-started/onboarding.md) for the complete workflow.

## Related Topics

- [Network Troubleshooting](network-issues.md)
- [Onboarding Issues](onboarding-issues.md)
- [Device Management](../devices/overview.md)
