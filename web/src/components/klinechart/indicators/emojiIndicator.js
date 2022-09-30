export const fruits = [
    '🍏', '🍎', '🍐', '🍊', '🍋', '🍌',
    '🍉', '🍇', '🍓', '🍈', '🍒', '🍑',
    '🍍', '🥥', '🥝', '🥭', '🥑', '🍏'
]

// custom indicator
export const emojiTechnicalIndicator = {
    name: 'EMOJI',
    plots: [
        { key: 'emoji' }
    ],
    calcTechnicalIndicator: (kLineDataList) => {
        const result = []
        kLineDataList.forEach(kLineData => {
            result.push({ emoji: kLineData.close, text: fruits[Math.floor(Math.random() * 17)] })
        })
        return result
    },
    render: ({ ctx, dataSource, viewport, xAxis, yAxis }) => {
        ctx.font = `${viewport.barSpace}px Helvetica Neue`
        ctx.textAlign = 'center'
        for (let i = dataSource.from; i < dataSource.to; i++) {
            const data = dataSource.technicalIndicatorDataList[i]
            const x = xAxis.convertToPixel(i)
            const y = yAxis.convertToPixel(data.emoji)
            ctx.fillText(data.text, x, y)
        }
    }
}