# iso-edit

The goal of this test program is to unpack a base iso and inject an additional cpio archive such that arbitrary files can be added to the initrd.

The repackaging step should, ideally do the same thing as:

```bash
genisoimage \
   -V rhcos-46.82.202010091720-0 \
   -c isolinux/boot.cat \
   -b isolinux/isolinux.bin \
   -no-emul-boot \
   -boot-load-size 4 \
   -boot-info-table \
   -eltorito-alt-boot \
   -e images/efiboot.img \
   -no-emul-boot \
   -o isos/test-iso.iso \
   /tmp/iso-test/
```

## How to use

1. Download the base ISO from https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.6/4.6.1/rhcos-4.6.1-x86_64-live.x86_64.iso to the `isos` directory
2. Build `make build`
3. Run `./build/iso-edit`

## Notes

1. Running currently leaves the unpacked iso at `/tmp/iso-test`
2. The default location for the iso is `isos/my-rhcos.iso`
3. The input iso, output path, and files to add are all configurable with flags
4. `make run` is a convienience which will remove the temp dir and the default output iso, rebuild, and run with the defaults.
