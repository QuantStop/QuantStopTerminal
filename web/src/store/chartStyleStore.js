import { reactive } from 'vue'

export const defaultTextColorDark = '#808080FF'
export const defaultGridColorDark = '#292929'
export const defaultAxisLineColorDark = '#333333'
export const defaultCrossTextBackgroundColorDark = '#373a40'
export const defaultChartBackgroundColorDark = '#141414'

export const defaultTextColorLight = '#76808F'
export const defaultGridColorLight = '#ededed'
export const defaultAxisLineColorLight = '#DDDDDD'
export const defaultCrossTextBackgroundColorLight = '#686d76'
export const defaultChartBackgroundColorLight = '#ffffff'

export const chartStyle = reactive({
    themes: {
        dark: {
            textColor: defaultTextColorDark,
            gridColor: defaultGridColorDark,
            axisLineColor: defaultAxisLineColorDark,
            crossTextBackgroundColor: defaultCrossTextBackgroundColorDark,
            chartBackgroundColor: defaultChartBackgroundColorDark,
        },
        light: {
            textColor: defaultTextColorLight,
            gridColor: defaultGridColorLight,
            axisLineColor: defaultAxisLineColorLight,
            crossTextBackgroundColor: defaultCrossTextBackgroundColorLight,
            chartBackgroundColor: defaultChartBackgroundColorLight,
        },
    },
    props: {
        grid: {
            show: true,
            horizontal: {
                show: true,
                size: 1,
                color: '#393939',
                // 'solid'|'dash'
                style: 'dash',
                dashValue: [2, 2]
            },
            vertical: {
                show: true,
                size: 1,
                color: '#393939',
                // 'solid'|'dash'
                style: 'dash',
                dashValue: [2, 2]
            }
        },
        candle: {
            margin: {
                top: 0.2,
                bottom: 0.1
            },
            // 'candle_solid'|'candle_stroke'|'candle_up_stroke'|'candle_down_stroke'|'ohlc'|'area'
            type: 'candle_solid',
            bar: {
                upColor: '#26A69A',
                downColor: '#EF5350',
                noChangeColor: '#888888'
            },
            area: {
                lineSize: 2,
                lineColor: '#2196F3',
                value: 'close',
                backgroundColor: [{
                    offset: 0,
                    color: 'rgba(33, 150, 243, 0.01)'
                }, {
                    offset: 1,
                    color: 'rgba(33, 150, 243, 0.2)'
                }]
            },
            priceMark: {
                show: true,
                high: {
                    show: true,
                    color: '#D9D9D9',
                    textMargin: 5,
                    textSize: 10,
                    textFamily: 'Segoe UI',
                    textWeight: 'normal'
                },
                low: {
                    show: true,
                    color: '#D9D9D9',
                    textMargin: 5,
                    textSize: 10,
                    textFamily: 'Segoe UI',
                    textWeight: 'normal',
                },
                last: {
                    show: true,
                    upColor: '#26A69A',
                    downColor: '#EF5350',
                    noChangeColor: '#888888',
                    line: {
                        show: true,
                        // 'solid'|'dash'
                        style: 'dash',
                        dashValue: [4, 4],
                        size: 1
                    },
                    text: {
                        show: true,
                        size: 12,
                        paddingLeft: 2,
                        paddingTop: 2,
                        paddingRight: 2,
                        paddingBottom: 2,
                        color: '#FFFFFF',
                        family: 'Segoe UI',
                        weight: 'normal',
                        borderRadius: 2
                    }
                }
            },
            tooltip: {
                // 'always' | 'follow_cross' | 'none'
                showRule: 'always',
                // 'standard' | 'rect'
                showType: 'standard',
                labels: ['T: ', 'O: ', 'C: ', 'H: ', 'L: ', 'V: '],
                values: null,
                defaultValue: 'n/a',
                rect: {
                    paddingLeft: 0,
                    paddingRight: 0,
                    paddingTop: 0,
                    paddingBottom: 6,
                    offsetLeft: 8,
                    offsetTop: 8,
                    offsetRight: 8,
                    borderRadius: 4,
                    borderSize: 1,
                    borderColor: '#3f4254',
                    backgroundColor: 'rgba(17, 17, 17, .3)'
                },
                text: {
                    size: 12,
                    family: 'Segoe UI',
                    weight: 'normal',
                    color: '#D9D9D9',
                    marginLeft: 8,
                    marginTop: 6,
                    marginRight: 8,
                    marginBottom: 0
                }
            }
        },
        technicalIndicator: {
            margin: {
                top: 0.2,
                bottom: 0.1
            },
            bar: {
                upColor: '#26A69A',
                downColor: '#EF5350',
                noChangeColor: '#888888'
            },
            line: {
                size: 1,
                colors: ['#FF9600', '#9D65C9', '#2196F3', '#E11D74', '#01C5C4']
            },
            circle: {
                upColor: '#26A69A',
                downColor: '#EF5350',
                noChangeColor: '#888888'
            },
            lastValueMark: {
                show: true,
                text: {
                    show: true,
                    color: '#ffffff',
                    size: 12,
                    family: 'Segoe UI',
                    weight: 'normal',
                    paddingLeft: 3,
                    paddingTop: 2,
                    paddingRight: 3,
                    paddingBottom: 2,
                    borderRadius: 2
                }
            },
            tooltip: {
                // 'always' | 'follow_cross' | 'none'
                showRule: 'always',
                // 'standard' | 'rect'
                showType: 'standard',
                showName: true,
                showParams: true,
                defaultValue: 'n/a',
                text: {
                    size: 12,
                    family: 'Segoe UI',
                    weight: '400',
                    color: '#D9D9D9',
                    marginTop: 6,
                    marginRight: 8,
                    marginBottom: 0,
                    marginLeft: 8
                }
            }
        },
        xAxis: {
            show: true,
            height: null,
            axisLine: {
                show: true,
                color: '#888888',
                size: 1
            },
            tickText: {
                show: true,
                color: '#808080FF',
                family: 'Segoe UI',
                weight: '400',
                lineHeight: 1.5,
                size: 12,
                paddingTop: 3,
                paddingBottom: 6
            },
            tickLine: {
                show: true,
                size: 1,
                length: 3,
                color: '#888888'
            }
        },
        yAxis: {
            show: true,
            width: null,
            // 'left' | 'right'
            position: 'right',
            // 'normal' | 'percentage' | 'log'
            type: 'normal',
            inside: false,
            axisLine: {
                show: true,
                color: '#888888',
                size: 1
            },
            tickText: {
                show: true,
                color: '#808080FF',
                family: 'Segoe UI',
                weight: '400',
                lineHeight: 1.5,
                size: 12,
                paddingLeft: 3,
                paddingRight: 6
            },
            tickLine: {
                show: true,
                size: 1,
                length: 3,
                color: '#888888'
            }
        },
        separator: {
            size: 1,
            color: '#888888',
            fill: true,
            activeBackgroundColor: 'rgba(230, 230, 230, .15)'
        },
        crosshair: {
            show: true,
            horizontal: {
                show: true,
                line: {
                    show: true,
                    // 'solid'|'dash'
                    style: 'dash',
                    dashValue: [4, 2],
                    size: 1,
                    color: '#888888'
                },
                text: {
                    show: true,
                    color: '#D9D9D9',
                    size: 12,
                    family: 'Segoe UI',
                    weight: 'normal',
                    paddingLeft: 2,
                    paddingRight: 2,
                    paddingTop: 2,
                    paddingBottom: 2,
                    borderSize: 1,
                    borderColor: '#505050',
                    borderRadius: 2,
                    backgroundColor: '#505050'
                }
            },
            vertical: {
                show: true,
                line: {
                    show: true,
                    // 'solid'|'dash'
                    style: 'dash',
                    dashValue: [4, 2],
                    size: 1,
                    color: '#888888'
                },
                text: {
                    show: true,
                    color: '#D9D9D9',
                    size: 12,
                    family: 'Segoe UI',
                    weight: 'normal',
                    paddingLeft: 2,
                    paddingRight: 2,
                    paddingTop: 2,
                    paddingBottom: 2,
                    borderSize: 1,
                    borderColor: '#505050',
                    borderRadius: 2,
                    backgroundColor: '#505050'
                }
            }
        },
        shape: {
            point: {
                backgroundColor: '#2196F3',
                borderColor: '#2196F3',
                borderSize: 1,
                radius: 4,
                activeBackgroundColor: '#2196F3',
                activeBorderColor: '#2196F3',
                activeBorderSize: 1,
                activeRadius: 6
            },
            line: {
                // 'solid'|'dash'
                style: 'solid',
                color: '#2196F3',
                size: 1,
                dashValue: [2, 2]
            },
            polygon: {
                // 'stroke'|'fill'
                style: 'stroke',
                stroke: {
                    // 'solid'|'dash'
                    style: 'solid',
                    size: 1,
                    color: '#2196F3',
                    dashValue: [2, 2]
                },
                fill: {
                    color: 'rgba(33, 150, 243, 0.1)'
                }
            },
            arc: {
                // 'stroke'|'fill'
                style: 'stroke',
                stroke: {
                    // 'solid'|'dash'
                    style: 'solid',
                    size: 1,
                    color: '#2196F3',
                    dashValue: [2, 2]
                },
                fill: {
                    color: '#2196F3'
                }
            },
            text: {
                style: 'fill',
                color: '#2196F3',
                size: 12,
                family: 'Segoe UI',
                weight: 'normal',
                offset: [0, 0]
            }
        },
        annotation: {
            // 'top' | 'bottom' | 'point'
            position: 'top',
            offset: [20, 0],
            symbol: {
                // 'diamond' | 'circle' | 'rect' | 'triangle' | 'custom' | 'none'
                type: 'diamond',
                size: 8,
                color: '#2196F3',
                activeSize: 10,
                activeColor: '#FF9600'
            }
        },
        tag: {
            // 'top' | 'bottom' | 'point'
            position: 'point',
            offset: 0,
            line: {
                show: true,
                style: 'dash',
                dashValue: [4, 2],
                size: 1,
                color: '#2196F3'
            },
            text: {
                color: '#FFFFFF',
                backgroundColor: '#2196F3',
                size: 12,
                family: 'Segoe UI',
                weight: 'normal',
                paddingLeft: 2,
                paddingRight: 2,
                paddingTop: 2,
                paddingBottom: 2,
                borderRadius: 2,
                borderSize: 1,
                borderColor: '#2196F3'
            },
            mark: {
                offset: 0,
                color: '#FFFFFF',
                backgroundColor: '#2196F3',
                size: 12,
                family: 'Segoe UI',
                weight: 'normal',
                paddingLeft: 2,
                paddingRight: 2,
                paddingTop: 2,
                paddingBottom: 2,
                borderRadius: 2,
                borderSize: 1,
                borderColor: '#2196F3'
            }
        },
    },
    setThemeOptions(theme) {

        let textColor = defaultTextColorLight
        let gridColor = defaultGridColorLight
        let axisLineColor = defaultAxisLineColorLight
        let crossLineColor = defaultAxisLineColorLight
        let crossTextBackgroundColor = defaultCrossTextBackgroundColorLight
        let chartBackgroundColor = defaultChartBackgroundColorLight

        switch (theme) {
            case "dark":
                textColor = this.themes.dark.textColor
                gridColor = this.themes.dark.gridColor
                axisLineColor = this.themes.dark.axisLineColor
                crossLineColor = this.themes.dark.axisLineColor
                crossTextBackgroundColor = this.themes.dark.crossTextBackgroundColor
                chartBackgroundColor = this.themes.dark.chartBackgroundColor
                break;
            case "light":
                textColor = this.themes.light.textColor
                gridColor = this.themes.light.gridColor
                axisLineColor = this.themes.light.axisLineColor
                crossLineColor = this.themes.light.axisLineColor
                crossTextBackgroundColor = this.themes.light.crossTextBackgroundColor
                chartBackgroundColor = this.themes.light.chartBackgroundColor
                break;
            default:
                break;
        }

        this.props.grid.horizontal.color = gridColor
        this.props.grid.vertical.color = gridColor
        this.props.candle.priceMark.high.color = textColor
        this.props.candle.priceMark.low.color = textColor
        this.props.candle.tooltip.text.color = textColor
        this.props.technicalIndicator.tooltip.text.color = textColor
        this.props.xAxis.axisLine.color = axisLineColor
        this.props.xAxis.tickLine.color = axisLineColor
        this.props.xAxis.tickText.color = textColor
        this.props.yAxis.axisLine.color = axisLineColor
        this.props.yAxis.tickLine.color = axisLineColor
        this.props.yAxis.tickText.color = textColor
        this.props.separator.color = axisLineColor
        this.props.crosshair.horizontal.line.color = crossLineColor
        this.props.crosshair.horizontal.text.backgroundColor = crossTextBackgroundColor
        this.props.crosshair.vertical.line.color = crossLineColor
        this.props.crosshair.vertical.text.backgroundColor = crossTextBackgroundColor

    },

    getStyle() {
        if (localStorage.getItem('chartStyle')) {
            try {
                let p = JSON.parse(localStorage.getItem('chartStyle'))
                this.props = p.props;
                this.themes = p.themes;
            } catch(e) {
                localStorage.removeItem('chartStyle');
            }
        } else {
            // first time no local storage saved yet, save the defaults
            this.saveStyle()
        }
        return this.props
    },

    saveStyle() {
        try {
            let p = {
                themes: this.themes,
                props: this.props
            }
            localStorage.setItem('chartStyle', JSON.stringify(p));
        } catch(e) {
            localStorage.removeItem('chartStyle');
        }
    },
})











