<!--suppress ALL -->
<div id="top"></div>

<!-- PROJECT SHIELDS -->
<!--
Using markdown "reference style" links for readability.
Reference links are enclosed in brackets [ ] instead of parentheses ( ).
See the bottom of this document for the declaration of the reference variables
for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
https://www.markdownguide.org/basic-syntax/#reference-style-links

Note: 
  To have badges be centered there MUST be a blank line between markdown and div tags
  https://stackoverflow.com/questions/70292850/centre-align-shield-io-in-github-readme-file
-->
<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

</div>

<!-- PROJECT LOGO -->
<div align="center">
  <h1>QuantstopExchange</h1>

  <p align="center">
    Go library for crypto currency and stock market exchanges. <br />
  </p>
</div>

<br />

<!-- ABOUT THE PROJECT -->
## About The Project

Building financial applications is hard. <br>
QuantstopExchange, abbreviated as qsx, provides a simple way to interact with many exchanges, and ensure a common data structure. 



<!-- GETTING STARTED -->
## Getting Started
Import the quantstopexchange factory, the core library, and any vendor exchanges you wish to use:
```go
import (
    "github.com/quantstop/quantstopexchange"
    "github.com/quantstop/quantstopterminal/pkg/quantstopexchange/qsx"
    "github.com/quantstop/quantstopterminal/pkg/quantstopexchange/vendors/coinbasepro"
)
```

Create an exchange object: <br>
<i><small>Authentication is only needed for operations that require it. <br>
You may leave the fields empty as blank input strings. <br>
</small></i>

```go
config := &qsx.Config{
    Auth:    qsx.NewAuth("key", "pass", "secret"),
    Sandbox: true,
}
coinbasepro, err := quantstopexchange.NewExchange(qsx.CoinbasePro, config)
if err != nil {
    // handle error
}
```

Use the exchange object to perform an action:

```go
candles, err := coinbasepro.GetHistoricalCandles(context.TODO(), "BTC-USD", "1m")
if err != nil {
    // handle error
}
for _, candle := range candles {
    fmt.Println(fmt.Sprintf("Close: %v", candle.Close))
}
```



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create.
Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request.
You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request




<!-- Eula_en-us.rtf -->
## License

Distributed under the MIT License. See the full [LICENSE](LICENSE) for more information.




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/quantstop/qsx.svg?style=for-the-badge
[contributors-url]: https://github.com/quantstop/qsx/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/quantstop/qsx.svg?style=for-the-badge
[forks-url]: https://github.com/quantstop/qsx/network/members
[stars-shield]: https://img.shields.io/github/stars/quantstop/qsx.svg?style=for-the-badge
[stars-url]: https://github.com/quantstop/qsx/stargazers
[issues-shield]: https://img.shields.io/github/issues/quantstop/qsx.svg?style=for-the-badge
[issues-url]: https://github.com/quantstop/qsx/issues
[license-shield]: https://img.shields.io/github/license/quantstop/qsx.svg?style=for-the-badge
[license-url]: https://github.com/quantstop/qsx/blob/main/LICENSE
