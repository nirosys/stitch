module github.com/nirosys/stitch

go 1.13

require (
	github.com/emicklei/dot v0.10.1
	github.com/gookit/color v1.2.5
	github.com/nirosys/gaufre v0.0.0-20200724171953-6f0402dac517
	github.com/peterh/liner v0.0.0-00010101000000-000000000000
	github.com/rs/xid v1.2.1
	github.com/spf13/cobra v0.0.7
)

replace github.com/peterh/liner => ./cmd/stitch/subcmd/internal/liner
