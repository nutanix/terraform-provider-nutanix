## 1.0.2 (Unreleased)
**Fixed bugs:**

- user\_data typos cause panic, ForceNew for guest\_customization fields [\#67](https://github.com/terraform-providers/terraform-provider-nutanix/issues/67)
- Example config is not valid for power\_state [\#9](https://github.com/terraform-providers/terraform-provider-nutanix/issues/9)

**Closed issues:**

- Do not shutdown machine on certain changes [\#74](https://github.com/terraform-providers/terraform-provider-nutanix/issues/74)
- Update compatibility matrix for TF provider [\#71](https://github.com/terraform-providers/terraform-provider-nutanix/issues/71)
- Terraform 0.12 compatibility [\#66](https://github.com/terraform-providers/terraform-provider-nutanix/issues/66)
- Unable to specify category key names dynamically in resource:nutanix\_virtual\_machine [\#61](https://github.com/terraform-providers/terraform-provider-nutanix/issues/61)
- data source nutanix\_virtual\_machine fails with "Invalid address to set: \[\]string{"nic\_list", "0", "is\_connected"}" [\#57](https://github.com/terraform-providers/terraform-provider-nutanix/issues/57)
- cannot unmarshal string into GO struct field MessageResource.details of the type map \[string\]interface{} [\#44](https://github.com/terraform-providers/terraform-provider-nutanix/issues/44)
- Error when re-applying a plan [\#41](https://github.com/terraform-providers/terraform-provider-nutanix/issues/41)
- unable to spin vm with error "'api\_version' is a required property" [\#36](https://github.com/terraform-providers/terraform-provider-nutanix/issues/36)
- Change VM caused disk being deleted [\#34](https://github.com/terraform-providers/terraform-provider-nutanix/issues/34)
- Unable to change VM resources when it has a network connected [\#33](https://github.com/terraform-providers/terraform-provider-nutanix/issues/33)
- Consider replacing satori/go.uuid [\#31](https://github.com/terraform-providers/terraform-provider-nutanix/issues/31)

**Merged pull requests:**

- Dynamic category name support [\#76](https://github.com/terraform-providers/terraform-provider-nutanix/pull/76) ([marinsalinas](https://github.com/marinsalinas))
- Patch to allow some changes to be hotplug [\#75](https://github.com/terraform-providers/terraform-provider-nutanix/pull/75) ([Jorge-Holgado](https://github.com/Jorge-Holgado))
- Terraform 0.12 Provider Support. [\#72](https://github.com/terraform-providers/terraform-provider-nutanix/pull/72) ([marinsalinas](https://github.com/marinsalinas))
- Fix typos causing panic on userdata change, ForceNew on change [\#68](https://github.com/terraform-providers/terraform-provider-nutanix/pull/68) ([rxacevedo](https://github.com/rxacevedo))
- Refactor: change update workflow to use spec instead status in virtual machine resource. [\#62](https://github.com/terraform-providers/terraform-provider-nutanix/pull/62) ([marinsalinas](https://github.com/marinsalinas))
- Added fix for Issue \#57 [\#60](https://github.com/terraform-providers/terraform-provider-nutanix/pull/60) ([chandru-tkc](https://github.com/chandru-tkc))
- added VMNic.IsConnected to preFillResUpdateRequest [\#59](https://github.com/terraform-providers/terraform-provider-nutanix/pull/59) ([switchboardOp](https://github.com/switchboardOp))


## 1.0.1 (May 01, 2019)

**Implemented enhancements:**

- Ability to resize a vdisk when deploying from image resource [\#23](https://github.com/terraform-providers/terraform-provider-nutanix/issues/23)
- nutanix\_image: long running image create completes successfully at 10 minutes, without getting success call from API [\#20](https://github.com/terraform-providers/terraform-provider-nutanix/issues/20)
- Datasources for categories [\#8](https://github.com/terraform-providers/terraform-provider-nutanix/issues/8)

**Fixed bugs:**

- nutanix\\_image: long running image create completes successfully at 10 minutes, without getting success call from API [\#20](https://github.com/terraform-providers/terraform-provider-nutanix/issues/20)
- nutanix\_virtual\_machine, ip\_endpoint\_list not correctly working [\#12](https://github.com/terraform-providers/terraform-provider-nutanix/issues/12)
- nutanix\_clusters data source makes Terraform crash [\#10](https://github.com/terraform-providers/terraform-provider-nutanix/issues/10)

**Closed issues:**

- Terraform error trying to use nutanix\_virtual\_machine data source [\#30](https://github.com/terraform-providers/terraform-provider-nutanix/issues/30)
- \[PROPOSAL\] Switch to Go Modules [\#29](https://github.com/terraform-providers/terraform-provider-nutanix/issues/29)
- hard code \(kind\) to be appropriate kind for each resource [\#27](https://github.com/terraform-providers/terraform-provider-nutanix/issues/27)
- Implement HTTP proxy capability at provider level [\#26](https://github.com/terraform-providers/terraform-provider-nutanix/issues/26)
- Ability to add a serial port to a VM [\#25](https://github.com/terraform-providers/terraform-provider-nutanix/issues/25)
- Add additional acceptance test configurations for nutanix resources and data sources [\#24](https://github.com/terraform-providers/terraform-provider-nutanix/issues/24)
- Example main.tf - guest\_customization\_cloud\_init does not work [\#21](https://github.com/terraform-providers/terraform-provider-nutanix/issues/21)
- Can't apply just a single resource "nutanix\_image" [\#18](https://github.com/terraform-providers/terraform-provider-nutanix/issues/18)

**Merged pull requests:**

- fix category key data source basic test [\#55](https://github.com/terraform-providers/terraform-provider-nutanix/pull/55) ([marinsalinas](https://github.com/marinsalinas))
- Vdisk update \#23 [\#54](https://github.com/terraform-providers/terraform-provider-nutanix/pull/54) ([marinsalinas](https://github.com/marinsalinas))
- Data Source for categories [\#51](https://github.com/terraform-providers/terraform-provider-nutanix/pull/51) ([marinsalinas](https://github.com/marinsalinas))
- chore: update cibuild make formula to build for any OS \(windows, linu… [\#50](https://github.com/terraform-providers/terraform-provider-nutanix/pull/50) ([marinsalinas](https://github.com/marinsalinas))
- chore: fix bad urls in readme.md file [\#48](https://github.com/terraform-providers/terraform-provider-nutanix/pull/48) ([marinsalinas](https://github.com/marinsalinas))
- chore: fix bad urls in readme.md file [\#47](https://github.com/terraform-providers/terraform-provider-nutanix/pull/47) ([marinsalinas](https://github.com/marinsalinas))
- Add HTTP Proxy capability [\#46](https://github.com/terraform-providers/terraform-provider-nutanix/pull/46) ([marinsalinas](https://github.com/marinsalinas))
- Add Serial Port support in VM resource and data source [\#45](https://github.com/terraform-providers/terraform-provider-nutanix/pull/45) ([marinsalinas](https://github.com/marinsalinas))
- Fix: nutanix\_guest\_tools attributes [\#43](https://github.com/terraform-providers/terraform-provider-nutanix/pull/43) ([marinsalinas](https://github.com/marinsalinas))
- Image wait timeout \#20 [\#42](https://github.com/terraform-providers/terraform-provider-nutanix/pull/42) ([marinsalinas](https://github.com/marinsalinas))
- refactor: add nic\_list\_status to separate the computed values [\#40](https://github.com/terraform-providers/terraform-provider-nutanix/pull/40) ([marinsalinas](https://github.com/marinsalinas))
- Remove kind as a argument in cluster\_reference and metadata. [\#39](https://github.com/terraform-providers/terraform-provider-nutanix/pull/39) ([marinsalinas](https://github.com/marinsalinas))
- Fix clusters Data Source makes tf crash \#10 [\#38](https://github.com/terraform-providers/terraform-provider-nutanix/pull/38) ([marinsalinas](https://github.com/marinsalinas))
- \[MODULES\] Switch to Go Modules [\#37](https://github.com/terraform-providers/terraform-provider-nutanix/pull/37) ([appilon](https://github.com/appilon))
- website: change guest\_customization\_cloud\_init argument reference [\#32](https://github.com/terraform-providers/terraform-provider-nutanix/pull/32) ([marinsalinas](https://github.com/marinsalinas))
- \[AUTOMATED\] Upgrade to Go 1.11 [\#16](https://github.com/terraform-providers/terraform-provider-nutanix/pull/16) ([appilon](https://github.com/appilon))
- Update docs to reflect removal of network security rule resource [\#14](https://github.com/terraform-providers/terraform-provider-nutanix/pull/14) ([alias-dev](https://github.com/alias-dev))
- correcting example [\#7](https://github.com/terraform-providers/terraform-provider-nutanix/pull/7) ([olljanat](https://github.com/olljanat))
- Fix Spell and style [\#3](https://github.com/terraform-providers/terraform-provider-nutanix/pull/3) ([ryujisnote](https://github.com/ryujisnote))
- fix cluster datasource\(s\) if no http/smtp credentials are configured [\#1](https://github.com/terraform-providers/terraform-provider-nutanix/pull/1) ([simonfuhrer](https://github.com/simonfuhrer))

## 1.0.0 (October 09, 2018)

**Implemented enhancements:**

- Configure Test Coverage via Code Climate [\#112](https://github.com/nutanix/terraform-provider-nutanix/issues/112)
- Add Cluster DS [\#68](https://github.com/nutanix/terraform-provider-nutanix/issues/68)

**Fixed bugs:**

- client/v3/v3\_service.go - clean up bunch of TODO's [\#127](https://github.com/nutanix/terraform-provider-nutanix/issues/127)
- Incorporate gofmt simplify \(gofmt -s\) [\#91](https://github.com/nutanix/terraform-provider-nutanix/issues/91)
- Fix Go Report Card issues with latest merge [\#90](https://github.com/nutanix/terraform-provider-nutanix/issues/90)

**Closed issues:**

- Get test coverage up \(at least over 50%\) [\#141](https://github.com/nutanix/terraform-provider-nutanix/issues/141)
- Crash during VM change [\#132](https://github.com/nutanix/terraform-provider-nutanix/issues/132)
- Terraform crash during vm create [\#119](https://github.com/nutanix/terraform-provider-nutanix/issues/119)
- Feature requests: max length for metadata / userdata in guest customization [\#115](https://github.com/nutanix/terraform-provider-nutanix/issues/115)
- Terraform crash during vm re-create [\#114](https://github.com/nutanix/terraform-provider-nutanix/issues/114)
- Terraform client fails to read cloud init values [\#111](https://github.com/nutanix/terraform-provider-nutanix/issues/111)
- Volume Groups DS [\#100](https://github.com/nutanix/terraform-provider-nutanix/issues/100)
- Volume Group DS [\#99](https://github.com/nutanix/terraform-provider-nutanix/issues/99)
- Volume Group RS [\#98](https://github.com/nutanix/terraform-provider-nutanix/issues/98)
- Add Subnets DS [\#97](https://github.com/nutanix/terraform-provider-nutanix/issues/97)
- Add Images DS [\#96](https://github.com/nutanix/terraform-provider-nutanix/issues/96)
- Configuration for the Terraform Nutanix Documentation [\#88](https://github.com/nutanix/terraform-provider-nutanix/issues/88)
- Fix "GoLint/Naming/PackageNames" issue in nutanix/config.go [\#87](https://github.com/nutanix/terraform-provider-nutanix/issues/87)
- Fix "SC2046" issue in scripts/errcheck.sh [\#86](https://github.com/nutanix/terraform-provider-nutanix/issues/86)
- Fix "SC2046" issue in scripts/gofmtcheck.sh [\#85](https://github.com/nutanix/terraform-provider-nutanix/issues/85)
- Fix "SC2006" issue in scripts/gofmtcheck.sh [\#84](https://github.com/nutanix/terraform-provider-nutanix/issues/84)
- Fix "SC2006" issue in scripts/changelog-links.sh [\#83](https://github.com/nutanix/terraform-provider-nutanix/issues/83)
- Fix "GoLint/Naming/PackageNames" issue in nutanix/config.go [\#82](https://github.com/nutanix/terraform-provider-nutanix/issues/82)
- Test ticket from Code Climate [\#81](https://github.com/nutanix/terraform-provider-nutanix/issues/81)
- Add Network Security Rules DS [\#80](https://github.com/nutanix/terraform-provider-nutanix/issues/80)
- Add Network Security Rule DS [\#75](https://github.com/nutanix/terraform-provider-nutanix/issues/75)
- Terraform crash with VM and cloud init [\#72](https://github.com/nutanix/terraform-provider-nutanix/issues/72)
- Make all test passing [\#71](https://github.com/nutanix/terraform-provider-nutanix/issues/71)
- Add Cluster RS [\#69](https://github.com/nutanix/terraform-provider-nutanix/issues/69)
- Get Image resource test passing [\#67](https://github.com/nutanix/terraform-provider-nutanix/issues/67)
- implement issue template [\#61](https://github.com/nutanix/terraform-provider-nutanix/issues/61)
- Add codecov.io code coverage reporting [\#57](https://github.com/nutanix/terraform-provider-nutanix/issues/57)
- Add VMS Data Source [\#51](https://github.com/nutanix/terraform-provider-nutanix/issues/51)
- Terraform crash when updating image [\#49](https://github.com/nutanix/terraform-provider-nutanix/issues/49)
- Unclear protocol / error message for image upload [\#48](https://github.com/nutanix/terraform-provider-nutanix/issues/48)
- Cleanup all go report card issues [\#47](https://github.com/nutanix/terraform-provider-nutanix/issues/47)

**Merged pull requests:**

- merge develop june 20 [\#149](https://github.com/nutanix/terraform-provider-nutanix/pull/149) ([JonKohler](https://github.com/JonKohler))
- Add test cases for client and V3  \#141 [\#148](https://github.com/nutanix/terraform-provider-nutanix/pull/148) ([marinsalinas](https://github.com/marinsalinas))
- develop to master nightly [\#147](https://github.com/nutanix/terraform-provider-nutanix/pull/147) ([JonKohler](https://github.com/JonKohler))
- WIP: Test cov \#141 [\#146](https://github.com/nutanix/terraform-provider-nutanix/pull/146) ([Crizstian](https://github.com/Crizstian))
- Change LICENSE to MPL2 per hashi [\#143](https://github.com/nutanix/terraform-provider-nutanix/pull/143) ([JonKohler](https://github.com/JonKohler))
- lint changes [\#142](https://github.com/nutanix/terraform-provider-nutanix/pull/142) ([JonKohler](https://github.com/JonKohler))
- squash lint issues and add unit tests via gotest [\#140](https://github.com/nutanix/terraform-provider-nutanix/pull/140) ([JonKohler](https://github.com/JonKohler))
- pulling in all develop work on travis/gnu [\#138](https://github.com/nutanix/terraform-provider-nutanix/pull/138) ([JonKohler](https://github.com/JonKohler))
- Master [\#137](https://github.com/nutanix/terraform-provider-nutanix/pull/137) ([JonKohler](https://github.com/JonKohler))
- Develop [\#136](https://github.com/nutanix/terraform-provider-nutanix/pull/136) ([JonKohler](https://github.com/JonKohler))
- Codecov \#112 Add Travis CI Code Coverage configuration [\#135](https://github.com/nutanix/terraform-provider-nutanix/pull/135) ([marinsalinas](https://github.com/marinsalinas))
- Image upload fail [\#134](https://github.com/nutanix/terraform-provider-nutanix/pull/134) ([thetonymaster](https://github.com/thetonymaster))
- remove more dead code [\#131](https://github.com/nutanix/terraform-provider-nutanix/pull/131) ([JonKohler](https://github.com/JonKohler))
- remove old code [\#130](https://github.com/nutanix/terraform-provider-nutanix/pull/130) ([JonKohler](https://github.com/JonKohler))
- one last line split fix [\#129](https://github.com/nutanix/terraform-provider-nutanix/pull/129) ([JonKohler](https://github.com/JonKohler))
- Pulling in test updates having challenges in GCP, but pass locally [\#128](https://github.com/nutanix/terraform-provider-nutanix/pull/128) ([JonKohler](https://github.com/JonKohler))
- squashing develop and master [\#126](https://github.com/nutanix/terraform-provider-nutanix/pull/126) ([JonKohler](https://github.com/JonKohler))
- fix quotes again [\#125](https://github.com/nutanix/terraform-provider-nutanix/pull/125) ([JonKohler](https://github.com/JonKohler))
- fix quotes [\#124](https://github.com/nutanix/terraform-provider-nutanix/pull/124) ([JonKohler](https://github.com/JonKohler))
- fix backticks [\#123](https://github.com/nutanix/terraform-provider-nutanix/pull/123) ([JonKohler](https://github.com/JonKohler))
- removing wait for ip and ip\_address from the vm ds and rs since not e… [\#121](https://github.com/nutanix/terraform-provider-nutanix/pull/121) ([Crizstian](https://github.com/Crizstian))
- updating example file and binaries to latest [\#120](https://github.com/nutanix/terraform-provider-nutanix/pull/120) ([Crizstian](https://github.com/Crizstian))
- Develop [\#118](https://github.com/nutanix/terraform-provider-nutanix/pull/118) ([Crizstian](https://github.com/Crizstian))
- Code refactor [\#117](https://github.com/nutanix/terraform-provider-nutanix/pull/117) ([marinsalinas](https://github.com/marinsalinas))
- Network sec rules refactor [\#116](https://github.com/nutanix/terraform-provider-nutanix/pull/116) ([marinsalinas](https://github.com/marinsalinas))
- Cloud init custom keys \#111 [\#113](https://github.com/nutanix/terraform-provider-nutanix/pull/113) ([marinsalinas](https://github.com/marinsalinas))
- Upload Image [\#110](https://github.com/nutanix/terraform-provider-nutanix/pull/110) ([thetonymaster](https://github.com/thetonymaster))
- Volume group rs 98 [\#109](https://github.com/nutanix/terraform-provider-nutanix/pull/109) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#108](https://github.com/nutanix/terraform-provider-nutanix/pull/108) ([Crizstian](https://github.com/Crizstian))
- refactorizing vm ds and rs [\#107](https://github.com/nutanix/terraform-provider-nutanix/pull/107) ([Crizstian](https://github.com/Crizstian))
- refactoring nsr ds and rs [\#106](https://github.com/nutanix/terraform-provider-nutanix/pull/106) ([Crizstian](https://github.com/Crizstian))
- Develop [\#105](https://github.com/nutanix/terraform-provider-nutanix/pull/105) ([Crizstian](https://github.com/Crizstian))
- Refactor [\#104](https://github.com/nutanix/terraform-provider-nutanix/pull/104) ([Crizstian](https://github.com/Crizstian))
- Develop [\#103](https://github.com/nutanix/terraform-provider-nutanix/pull/103) ([Crizstian](https://github.com/Crizstian))
- 97 subnets [\#102](https://github.com/nutanix/terraform-provider-nutanix/pull/102) ([Crizstian](https://github.com/Crizstian))
- 96 image ds [\#101](https://github.com/nutanix/terraform-provider-nutanix/pull/101) ([Crizstian](https://github.com/Crizstian))
- Documentation [\#95](https://github.com/nutanix/terraform-provider-nutanix/pull/95) ([marinsalinas](https://github.com/marinsalinas))
- Documentation [\#94](https://github.com/nutanix/terraform-provider-nutanix/pull/94) ([marinsalinas](https://github.com/marinsalinas))
- Refactor: client check errors [\#93](https://github.com/nutanix/terraform-provider-nutanix/pull/93) ([marinsalinas](https://github.com/marinsalinas))
- Network security rules \#80 [\#89](https://github.com/nutanix/terraform-provider-nutanix/pull/89) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#79](https://github.com/nutanix/terraform-provider-nutanix/pull/79) ([Crizstian](https://github.com/Crizstian))
- 68 cluster ds [\#78](https://github.com/nutanix/terraform-provider-nutanix/pull/78) ([Crizstian](https://github.com/Crizstian))
- Network Security Rule Data Source \#75 [\#77](https://github.com/nutanix/terraform-provider-nutanix/pull/77) ([marinsalinas](https://github.com/marinsalinas))
- Assign categories \#40 [\#76](https://github.com/nutanix/terraform-provider-nutanix/pull/76) ([marinsalinas](https://github.com/marinsalinas))
- Test passing \#71 [\#74](https://github.com/nutanix/terraform-provider-nutanix/pull/74) ([marinsalinas](https://github.com/marinsalinas))
- fix issue with vm.GuestCustomization not being initialized [\#73](https://github.com/nutanix/terraform-provider-nutanix/pull/73) ([htj](https://github.com/htj))
- Image test \#67 [\#70](https://github.com/nutanix/terraform-provider-nutanix/pull/70) ([marinsalinas](https://github.com/marinsalinas))
- fixes image rs [\#66](https://github.com/nutanix/terraform-provider-nutanix/pull/66) ([Crizstian](https://github.com/Crizstian))
- Categories resource \#36 [\#65](https://github.com/nutanix/terraform-provider-nutanix/pull/65) ([marinsalinas](https://github.com/marinsalinas))
- add completion callback for http request [\#64](https://github.com/nutanix/terraform-provider-nutanix/pull/64) ([marinsalinas](https://github.com/marinsalinas))
- Network Security Rule Resource \#30 [\#63](https://github.com/nutanix/terraform-provider-nutanix/pull/63) ([marinsalinas](https://github.com/marinsalinas))
- Add configuration for CI [\#58](https://github.com/nutanix/terraform-provider-nutanix/pull/58) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#55](https://github.com/nutanix/terraform-provider-nutanix/pull/55) ([JonKohler](https://github.com/JonKohler))
- Go fmt code [\#54](https://github.com/nutanix/terraform-provider-nutanix/pull/54) ([gliptak](https://github.com/gliptak))
- Develop [\#53](https://github.com/nutanix/terraform-provider-nutanix/pull/53) ([Crizstian](https://github.com/Crizstian))
- adding clusters data source to the provider, example.tf modified usin… [\#52](https://github.com/nutanix/terraform-provider-nutanix/pull/52) ([Crizstian](https://github.com/Crizstian))
- Go report issues \#47 [\#50](https://github.com/nutanix/terraform-provider-nutanix/pull/50) ([marinsalinas](https://github.com/marinsalinas))
- Develop [\#46](https://github.com/nutanix/terraform-provider-nutanix/pull/46) ([Crizstian](https://github.com/Crizstian))
- 25 cluster ds [\#45](https://github.com/nutanix/terraform-provider-nutanix/pull/45) ([Crizstian](https://github.com/Crizstian))
- Cleanup golint redundant nil check [\#44](https://github.com/nutanix/terraform-provider-nutanix/pull/44) ([gliptak](https://github.com/gliptak))
- Add GoReportCard badge [\#43](https://github.com/nutanix/terraform-provider-nutanix/pull/43) ([gliptak](https://github.com/gliptak))
- Add instructions for building project [\#42](https://github.com/nutanix/terraform-provider-nutanix/pull/42) ([marinsalinas](https://github.com/marinsalinas))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*
