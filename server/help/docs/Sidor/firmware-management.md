# Firmware Management

Manage and update device firmware versions centrally.

## Available Firmware

The system automatically scans the filesystem for available firmware versions organized by device model:

- `provisioning/rutx/` - RUTX series firmware
- `provisioning/rutm/` - RUTM series firmware
- `provisioning/rut2/` - RUT2M series firmware
- `provisioning/rut9/` - RUT9M series firmware

## Firmware Updates

### Manual Update

1. Navigate to **Devices** page
2. Select device to update
3. Click **Update Firmware**
4. Choose target version
5. Confirm update

### During Onboarding

Firmware updates can be performed automatically during the [onboarding process](../getting-started/onboarding.md).

## Firmware File Format

Firmware files should use the naming convention:
```
RUTX_R_00.07.10.2_WEBUI.bin
```

## VXLAN Support

Some configurations require VXLAN kernel modules. The system automatically checks and installs VXLAN support when needed.

## See Also

- [Device Onboarding](../getting-started/onboarding.md)
- [Troubleshooting Firmware Updates](../troubleshooting/firmware-issues.md)
