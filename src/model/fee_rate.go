package model

type FeeRate struct {
    FastestFee  int `json:"fastestFee"`
    HalfHourFee int `json:"halfHourFee"`
    HourFee     int `json:"hourFee"`
    EconomyFee  int `json:"economyFee"`
    MinimumFee  int `json:"minimumFee"`
}
