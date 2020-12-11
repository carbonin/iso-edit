# iso-edit

The goal of this test program is to unpack a base iso and inject an additional cpio archive such that arbitrary files can be added to the initrd.

The repackaging step now runs the following commands, but ideally it should use the go-diskfs library to do this

```bash
genisoimage \
   -verbose \
   -V rhcos-46.82.202010091720-0 -volset rhcos-46.82.202010091720-0 \
   -rational-rock -J -joliet-long \
   -eltorito-boot isolinux/isolinux.bin \
   -eltorito-catalog isolinux/isolinux.cat \
   -no-emul-boot -boot-load-size 4 -boot-info-table \
   -eltorito-alt-boot \
   -efi-boot images/efiboot.img \
   -no-emul-boot \
   -o isos/test-iso.iso \
   /tmp/iso-test/

isohybrid isos/test-iso.iso
```

Also this?
https://github.com/coreos/coreos-assembler/blob/510dbec7b84aa45a646079fe3341e2b5925c0774/src/cmd-buildextend-live#L481-L511

## How to use

1. Download the base ISO from https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.6/4.6.1/rhcos-4.6.1-x86_64-live.x86_64.iso to the `isos` directory
2. Build `make build`
3. Run `./build/iso-edit`

## Notes

1. Running currently leaves the unpacked iso at `isos/my-rhcos`
2. The default location for the iso is `isos/my-rhcos.iso`
4. `make run` is a convienience which will remove the my-rhcos dir and the default output iso, rebuild, and run with the defaults.

