package cmd

import (
	"fmt"
	"Seeyoner/core"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "漏洞检测",
	Long: `漏洞检测功能
`,
	Run: func(cmd *cobra.Command, args []string) {
		factory := new(core.IFactory)
		if vulnId == 0 {
			for i :=1 ; i < 13; i++ {
				fmt.Print("[", i, "] >>> ")
				iScan := factory.NewFactory(i)
				iScan.Scan(url)
			}
		} else {
			iScan := factory.NewFactory(vulnId)
			iScan.Scan(url)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&url, "targetUrl", "u", "", "targetUrl")
	scanCmd.Flags().IntVarP(&vulnId, "vulnId", "i", 0, "vulnId")
	scanCmd.MarkFlagRequired("targetUrl")
	scanCmd.MarkFlagRequired("vulnId")
}
