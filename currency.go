package miser

func Currency(isoCode string) (bool, string, rune) {

	switch isoCode {
	case "USD":
		return true, "United States Dollar", '$'
	case "EUR":
		return true, "Euro", '€'
	case "GBP":
		return true, "United Kingdom Pound", '£'
	case "CUP":
		return true, "Cuba Peso", '₱'
	case "CNY":
		return true, "China Yuan Renminbi", '¥'
	case "JPY":
		return true, "Japan Yen", '¥'
	case "AZN":
		return true, "Azerbaijan Manat", '₼'
	case "TRY":
		return true, "Turkish Lira", '₺'
	case "RUB":
		return true, "Russian Ruble", '₽'
	case "KPW":
		return true, "Korea Won", '₩'
	case "LAK":
		return true, "Laos Kip", '₭'
	case "NGN":
		return true, "Nigeria Naira", '₦'
	case "THB":
		return true, "Thailand Baht", '฿'
	case "UAH":
		return true, "Ukraine Hryvnia", '₴'
	case "KZT":
		return true, "Kazakhstan Tenge", '₸'
	case "VND":
		return true, "Viet Nam Dong", '₫'
	case "AFN":
		return true, "Afghanistan Afghani", '؋'
	case "BDT":
		return true, "Bangladeshi taka", '৳'
	case "INR":
		return true, "Indian rupee", '₹'

	// other dollars:
	case "AUD":
		return true, "Australia Dollar", '$'
	case "BSD":
		return true, "Bahamas Dollar", '$'
	case "BBD":
		return true, "Barbados Dollar", '$'
	case "BZD":
		return true, "Belize Dollar", '$'
	case "BMD":
		return true, "Bermuda Dollar", '$'
	case "BND":
		return true, "Brunei Darussalam Dollar", '$'
	case "CAD":
		return true, "Canada Dollar", '$'
	case "KYD":
		return true, "Cayman Islands Dollar", '$'
	case "XCD":
		return true, "East Caribbean Dollar", '$'
	case "FJD":
		return true, "Fiji Dollar", '$'
	case "GYD":
		return true, "Guyana Dollar", '$'
	case "HKD":
		return true, "Hong Kong Dollar", '$'
	case "JMD":
		return true, "Jamaica Dollar", '$'
	case "LRD":
		return true, "Liberia Dollar", '$'
	case "NAD":
		return true, "Namibia Dollar", '$'
	case "NZD":
		return true, "New Zealand Dollar", '$'
	case "SGD":
		return true, "Singapore Dollar", '$'
	case "SBD":
		return true, "Solomon Islands", '$'
	case "SRD":
		return true, "Suriname Dollar", '$'
	case "TWD":
		return true, "Taiwan New Dollar", '$'
	case "TTD":
		return true, "Trinidad and Tobago Dollar", '$'
	case "TVD":
		return true, "Tuvalu Dollar", '$'
	case "ZWD":
		return true, "Zimbabwe Dollar", '$'

	// crypto currencies:
	case "BTC":
		return true, "Bitcoin", '₿'
	case "ETH":
		return true, "Ethereum", '⟠'
	case "USDT":
		return true, "Tether USDt", '₮'
	case "BNB":
		return true, "BNB", 0
	case "SOL":
		return true, "Solana", '◎'
	case "XRP":
		return true, "XRP", '✕'
	case "USDC":
		return true, "USDC", 0
	case "DOGE":
		return true, "Dogecoin", 'Ð'
	case "TON":
		return true, "Toncoin", 0
	case "ADA":
		return true, "Cardano", '₳'
	case "AVAX":
		return true, "Avalanche", 0
	case "TRX":
		return true, "TRON", 0
	case "BCH":
		return true, "Bitcoin Cash", 'Ƀ'
	case "XMR":
		return true, "Monero", 'ɱ'

	default:
		return false, "", 0

	}
}
