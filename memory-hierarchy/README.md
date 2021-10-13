
age

- [baseline](https://github.com/savarin/computer-systems/blob/5bf88d681aa56e0f8650c00756ae18dd68a9dd0f/memory-hierarchy/metrics.go)	2207043 ns/op
- [convert to array](https://github.com/savarin/computer-systems/commit/d25677e18a31621262e6c151a56079964f45d6cc?branch=d25677e18a31621262e6c151a56079964f45d6cc&diff=split)	535163 ns/op
- [non-overflow safe sum](https://github.com/savarin/computer-systems/commit/136c2af1bf958c12f26544047f4c8ac5c9877ef3?branch=136c2af1bf958c12f26544047f4c8ac5c9877ef3&diff=split)	37980 ns/op
- [iteration with index](https://github.com/savarin/computer-systems/commit/78dea3280bbccbd3da61c3c7f1baf1d02a89e1f3?branch=78dea3280bbccbd3da61c3c7f1baf1d02a89e1f3&diff=split)	36071 ns/op
- [use uint8](https://github.com/savarin/computer-systems/commit/bbff04d1aaf6e545170384dedbf37ac6a4529872?branch=bbff04d1aaf6e545170384dedbf37ac6a4529872&diff=split)	30283 ns/op
- [loop unrolling](https://github.com/savarin/computer-systems/commit/b181e048307c5b6a947bdde545ab8e8d11ccb4ff?branch=b181e048307c5b6a947bdde545ab8e8d11ccb4ff&diff=split)	27246 ns/op
- [remove bound checks](https://github.com/savarin/computer-systems/commit/879aae6fe0ab4d2edba7cc37133e8498c29ae034?branch=879aae6fe0ab4d2edba7cc37133e8498c29ae034&diff=split)	18240 ns/op


payments

- [baseline](https://github.com/savarin/computer-systems/blob/097cddffbcdb2d3e4b7c20579a75335f4227ea3c/memory-hierarchy/metrics.go) 29184902 / 56645657 ns/op
- [convert to array](https://github.com/savarin/computer-systems/commit/809b74f086014e593317039d704a3971e5545154?branch=809b74f086014e593317039d704a3971e5545154&diff=split)	5398410 / 6379861 ns/op
- [non-overflow safe sum](https://github.com/savarin/computer-systems/commit/8cfbe02f22d8ecd3e8ccb9bd3f1002bfde1967dc?branch=8cfbe02f22d8ecd3e8ccb9bd3f1002bfde1967dc&diff=split)	995995 / 1997440 ns/op
- [use float32](https://github.com/savarin/computer-systems/commit/adaa6497c5573f03cd882d15b62a4ed7d054572a?branch=adaa6497c5573f03cd882d15b62a4ed7d054572a&diff=split)	976023 / 1957138 ns/op
