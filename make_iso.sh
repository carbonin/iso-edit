#!/bin/bash

pushd /tmp/iso-test

genisoimage \
   -V rhcos-46.82.202009222340-0 \
   -c isolinux/boot.cat \
   -b isolinux/isolinux.bin \
   -no-emul-boot \
   -boot-load-size 4 \
   -boot-info-table \
   -eltorito-alt-boot \
   -e images/efiboot.img \
   -no-emul-boot \
   -o ~/Documents/scratch/golang/iso-edit/isos/test-iso.iso \
   /tmp/iso-test/

popd
