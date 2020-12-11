This is a test branch to sort out what looks like a bug in go-diskfs

Currently some files extracted from a coreos live iso are only a fraction (typically about half) of the size they're supposed to be.

Download the base ISO from https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/4.6/4.6.1/rhcos-4.6.1-x86_64-live.x86_64.iso to the isos directory
Running `make run` in this branch will unpack isos/rhcos-4.6.1-x86_64-live.x86_64.iso into isos/my-rhcos.
Additionally lots of debugging info will print showing the difference between the size of the file written and what the file info struct shows.
As of this writing, that output looks something like this:

```
rm -rf build/*
rm -rf isos/my-rhcos*
go build -o build/iso-edit cmd/main.go
./build/iso-edit
Disk: &{File:0xc000010030 Info:0xc00007ec30 Type:0 Size:918552576 LogicalBlocksize:512 PhysicalBlocksize:512 Table:0xc00007c720 Writable:false DefaultBlocks:true}
FileSystem: &{workspace: size:918552576 start:0 file:0xc000010030 blocksize:2048 volumes:{descriptors:[<nil> <nil> 0xc00015e000 0xc00001c4dc 0xc000162000] primary:0xc00015e000} pathTable:0xc00000e3e0 rootDir:0xc000114070 suspEnabled:true suspSkip:0 suspExtensions:[0xc00007a190]}

Opening file: /EFI/redhat/grub.cfg
wrote 1061 bytes to isos/my-rhcos/EFI/redhat/grub.cfg, size: 1061

Opening file: /images/efiboot.img
wrote 4128768 bytes to isos/my-rhcos/images/efiboot.img, size: 8173158

Opening file: /images/ignition.img
wrote 163840 bytes to isos/my-rhcos/images/ignition.img, size: 262144

Opening file: /images/pxeboot/initrd.img
wrote 40140800 bytes to isos/my-rhcos/images/pxeboot/initrd.img, size: 80239428

Opening file: /images/pxeboot/rootfs.img
wrote 410025984 bytes to isos/my-rhcos/images/pxeboot/rootfs.img, size: 819969024

Opening file: /images/pxeboot/vmlinuz
wrote 4489216 bytes to isos/my-rhcos/images/pxeboot/vmlinuz, size: 8924528

Opening file: /isolinux/boot.cat
wrote 2048 bytes to isos/my-rhcos/isolinux/boot.cat, size: 2048

Opening file: /isolinux/boot.msg
wrote 58 bytes to isos/my-rhcos/isolinux/boot.msg, size: 58

Opening file: /isolinux/isolinux.bin
wrote 38912 bytes to isos/my-rhcos/isolinux/isolinux.bin, size: 38912

Opening file: /isolinux/isolinux.cfg
wrote 2047 bytes to isos/my-rhcos/isolinux/isolinux.cfg, size: 2047

Opening file: /isolinux/ldlinux.c32
wrote 98304 bytes to isos/my-rhcos/isolinux/ldlinux.c32, size: 116320

Opening file: /isolinux/libcom32.c32
wrote 131072 bytes to isos/my-rhcos/isolinux/libcom32.c32, size: 180476

Opening file: /isolinux/libutil.c32
wrote 23924 bytes to isos/my-rhcos/isolinux/libutil.c32, size: 23924

Opening file: /isolinux/vesamenu.c32
wrote 26724 bytes to isos/my-rhcos/isolinux/vesamenu.c32, size: 26724

Opening file: /zipl.prm
wrote 132 bytes to isos/my-rhcos/zipl.prm, size: 132
```
