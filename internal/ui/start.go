package ui

import (
	c "github.com/achannarasappa/ticker/v5/internal/common"
	mon "github.com/achannarasappa/ticker/v5/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
)

// Start launches the command line interface and starts capturing input
func Start(dep *c.Dependencies, ctx *c.Context) func() error {
	return func() error {

		monitors, _ := mon.NewMonitor(mon.ConfigMonitor{
			RefreshInterval: ctx.Config.RefreshInterval,
			TargetCurrency:  ctx.Config.Currency,
			Logger:          ctx.Logger,
			ConfigMonitorsYahoo: mon.ConfigMonitorsYahoo{
				BaseURL:           dep.MonitorYahooBaseURL,
				SessionRootURL:    dep.MonitorYahooSessionRootURL,
				SessionCrumbURL:   dep.MonitorYahooSessionCrumbURL,
				SessionConsentURL: dep.MonitorYahooSessionConsentURL,
			},
			ConfigMonitorPriceCoinbase: mon.ConfigMonitorPriceCoinbase{
				BaseURL:      dep.MonitorPriceCoinbaseBaseURL,
				StreamingURL: dep.MonitorPriceCoinbaseStreamingURL,
			},
		})

		p := tea.NewProgram(
			NewModel(*dep, *ctx, monitors),
			tea.WithMouseCellMotion(),
			tea.WithAltScreen(),
		)

		var err error

		err = monitors.SetOnUpdate(mon.ConfigUpdateFns{
			OnUpdateAssetQuote: func(symbol string, assetQuote c.AssetQuote, versionVector int) {
				p.Send(SetAssetQuoteMsg{
					symbol:        symbol,
					assetQuote:    assetQuote,
					versionVector: versionVector,
				})
			},
			OnUpdateAssetGroupQuote: func(assetGroupQuote c.AssetGroupQuote, versionVector int) {
				p.Send(SetAssetGroupQuoteMsg{
					assetGroupQuote: assetGroupQuote,
					versionVector:   versionVector,
				})
			},
		})

		if err != nil {

			return err
		}

		_, err = p.Run()

		return err
	}

}
