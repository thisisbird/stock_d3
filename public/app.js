var margin = { top: 20, right: 50, bottom: 30, left: 60 },
        width = 1300 - margin.left - margin.right,
        height = 500 - margin.top - margin.bottom;

// 設定時間格式
var parseDate = d3.timeParse("%Y/%m/%d");
// var parseDate = d3.timeParse("%Y%m%d");
var monthDate = d3.timeParse("%Y%m");



// K線圖的x
var x = techan.scale.financetime()
        .range([0, width]);
var crosshairY = d3.scaleLinear()
        .range([height, 0]);
// K線圖的y
var y = d3.scaleLinear()
        .range([height - 60, 0]);
// 成交量的y
var yVolume = d3.scaleLinear()
        .range([height, height - 60]);
//成交量的x
var xScale = d3.scaleBand().range([0, width]).padding(0.15);

var tradearrow = techan.plot.tradearrow()
        .xScale(x)
        .yScale(y);//箭頭指標
// .y(function(d) {//用了zoomin的高度會鎖死
//     // Display the buy and sell arrows a bit above and below the price, so the price is still visible
//     if(d.type === 'buy') return y(d.low)+5;
//     if(d.type === 'sell') return y(d.high)-5;
//     else return y(d.price);
// });

var smac = techan.plot.sma()
        .xScale(x)
        .yScale(y);

var sma0 = techan.plot.sma()
        .xScale(x)
        .yScale(y);

var sma1 = techan.plot.sma()
        .xScale(x)
        .yScale(y);
var ema2 = techan.plot.ema()
        .xScale(x)
        .yScale(y);
var candlestick = techan.plot.candlestick()
        .xScale(x)
        .yScale(y);

var zoom = d3.zoom()
        .scaleExtent([1, 20]) //設定縮放大小1 ~ 20倍
        .translateExtent([[0, 0], [width, height]]) // 設定可以縮放的範圍，註解掉就可以任意拖曳
        .extent([[margin.left, margin.top], [width, height]])
        .on("zoom", zoomed);

var zoomableInit, yInit;
var xAxis = d3.axisBottom()
        .scale(x);

var yAxis = d3.axisLeft()
        .scale(y);

var volumeAxis = d3.axisLeft(yVolume)
        .ticks(4)
        .tickFormat(d3.format(",.3s"));
var ohlcAnnotation = techan.plot.axisannotation()
        .axis(yAxis)
        .orient('left')
        .format(d3.format(',.2f'));
var timeAnnotation = techan.plot.axisannotation()
        .axis(xAxis)
        .orient('bottom')
        .format(d3.timeFormat('%Y-%m-%d'))
        .translate([0, height]);

// 設定十字線
var crosshair = techan.plot.crosshair()
        .xScale(x)
        .yScale(crosshairY)
        .xAnnotation(timeAnnotation)
        .yAnnotation(ohlcAnnotation)
        .on("move", move);

// 設定文字區域
var textSvg = d3.select("body").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");
//設定顯示文字，web版滑鼠拖曳就會顯示，App上則是要點擊才會顯示
var svgText = textSvg.append("g")
        .attr("class", "description")
        .append("text")
        //     .attr("x", margin.left)
        .attr("y", 6)
        .attr("dy", ".71em")
        .style("text-anchor", "start")
        .text("");
//設定畫圖區域
var svg = d3.select("body")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .attr("pointer-events", "all")
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");


var dataArr;
var dataMaArr_o;
var dataBuySellArr_o;
loadJSON("new.csv", "ma.csv","buyAndSell.csv","date");

function loadJSON(file,file2,file3 , type) {
        var date_start = $('#start').val();
        var date_end = $('#end').val();
        date_start = monthDate(date_start)
        date_end = monthDate(date_end)
        date_end.setMonth(date_end.getMonth() + 1)
        svg.selectAll("*").remove(); // 切換不同資料需要重新畫圖，因此需要先清除原先的圖案
        d3.csv(file2, function (error, data) {
                dataMaArr_o = data;//匯入ma資料
        });
        d3.csv(file3, function (error, data) {
                dataBuySellArr_o = data;//匯入買賣資料
        });
        d3.csv(file, function (error, data) {
                var accessor = candlestick.accessor();
                var jsonData = data;

                data =
                        jsonData
                                .map(function (d) { // 設定data的格式
                                        if (type == "date") {
                                                return {
                                                        date: parseDate(d["交易日期"]),
                                                        open: +d["開盤價"],
                                                        high: +d["最高價"],
                                                        low: +d["最低價"],
                                                        close: +d["收盤價"],
                                                        volume: +d["成交量"],
                                                        change: +d["漲跌價"],
                                                        percentChange: +d["漲跌%"].replace('%', ''),
                                                };
                                        } else {
                                                return {
                                                        date: monthDate(d[0]),
                                                        open: +d[3],
                                                        high: +d[4],
                                                        low: +d[5],
                                                        close: +d[6],
                                                        volume: +d[10],
                                                        change: +d[7],
                                                        percentChange: +d[8],
                                                };
                                        }

                                }).sort(function (a, b) { return d3.ascending(accessor.d(a), accessor.d(b)); });



                var newData = jsonData.map(function (d) {
                        if (type == "date") {
                                return {
                                        date: parseDate(d["交易日期"]),
                                        volume: d["成交量"]
                                }
                        } else {
                                return {
                                        date: monthDate(d[0]),
                                        volume: d[10]
                                }
                        }

                });
                //     .reverse();

               

                svg.append("g")
                        .attr("class", "candlestick");
                svg.append("g")
                        .attr("class", "sma ma-c get");
                svg.append("g")
                        .attr("class", "sma ma-c loss");
                // svg.append("g")
                //         .attr("class", "sma ma-0");
                // svg.append("g")
                //         .attr("class", "sma ma-1");
                // svg.append("g")
                //         .attr("class", "ema ma-2");
                svg.append("g")
                        .attr("class", "volume axis");
                svg.append("g")
                        .attr("class", "x axis")
                        .attr("transform", "translate(0," + height + ")");

                svg.append("g")
                        .attr("class", "y axis")
                        .append("text")
                        .attr("y", -10)
                        .style("text-anchor", "end")
                        .text("Price (TWD)");
                svg.append("g")
                        .attr("class", "tradearrow");
                // Data to display initially

                // start = 0
                // end = data.length;
                data = data.filter(function(value) {//k棒
                        return  date_start <= value.date && value.date <= date_end ;
                      });
                newData = newData.filter(function(value) {//成交量
                        return  date_start <= value.date && value.date <= date_end ;
                });
                
                dataMaArr = dataMaArr_o.filter(function(value) {//MA數字資料
                        return  date_start <= parseDate(value['交易日期']) && parseDate(value['交易日期']) <= date_end ;
                });
                
                dataBuySellArr = dataBuySellArr_o.filter(function(value) {//MA數字資料
                        return  date_start <= parseDate(value['交易日期']) && parseDate(value['交易日期']) <= date_end ;
                });
                // draw(data.slice(start, end), newData.slice(start, end));
                draw(data, newData);

        });
}



function draw(data, volumeData) {
        // 設定domain，決定各座標所用到的資料
        x.domain(data.map(candlestick.accessor().d));
        y.domain(techan.scale.plot.ohlc(data, candlestick.accessor()).domain());
        xScale.domain(volumeData.map(function (d) { return d.date; }))
        yVolume.domain(techan.scale.plot.volume(data).domain());

        var testData = [
                // {date : parseDate('2020/11/07'),value:10000,type:'buy'},
                // {date : parseDate('2020/11/09'),value:12000,type:'buy'},
                // {date : null,value:null},
                // {date : parseDate('2020/12/03'),value:12000,type:'sell'},
                // {date : parseDate('2020/12/01'),value:10000,type:'sell'},
                // {date : null,value:null},
                // {date : parseDate('2020/01/09'),value:14000,type:'sell'},
                // {date : parseDate('2020/01/21'),value:11000,type:'sell'},
                // {date : null,value:null},

        ]

        


        var trades = dataBuySellArr.map(function (d) {
                return {
                        date: parseDate(d["交易日期"]),
                        type: d["買賣"].trim(),
                        price: parseInt(d["價格"]),
                }
        });
        var get_line = [];
        var loss_line = [];
        var temp = [];
        dataBuySellArr.forEach(function (d) {
                temp.push(d)
                if(d["收益"] > 0){
                        temp.push({date:null,value:0});
                        get_line = get_line.concat(temp);
                        temp = [];
                }
                if(d["收益"] < 0){
                        temp.push({date:null,value:0});
                        loss_line = loss_line.concat(temp);
                        temp = []
                }
        });
        
        var get_line = get_line.map(function (d) {
                return {
                        date: parseDate(d["交易日期"]) ?? null,
                        value: d["價格"] === undefined ? null: parseInt(d["價格"])
                }
        });
        var loss_line = loss_line.map(function (d) {
                return {
                        date: parseDate(d["交易日期"]) ?? null,
                        value: d["價格"] === undefined ? null: parseInt(d["價格"]),
                }
        });
        console.log(get_line);
        console.log(loss_line);
        svg.select("g.sma.ma-c.get").attr("clip-path", "url(#candlestickClip)").datum(get_line).call(sma0);//賺線
        svg.select("g.sma.ma-c.loss").attr("clip-path", "url(#candlestickClip)").datum(loss_line).call(sma0);//賠線


        // Add a clipPath: everything out of this area won't be drawn.
        var clip = svg.append("defs").append("svg:clipPath")
                .attr("id", "clip")
                .append("svg:rect")
                .attr("width", width)
                .attr("height", height)
                .attr("x", 0)
                .attr("y", 0);
        // 針對K線圖的，讓他不會蓋到成交量bar chart
        var candlestickClip = svg.append("defs").append("svg:clipPath")
                .attr("id", "candlestickClip")
                .append("svg:rect")
                .attr("width", width)
                .attr("height", height - 60)
                .attr("x", 0)
                .attr("y", 0);

        xScale.range([0, width].map(d => d)); // 設定xScale回到初始值
        var chart = svg.selectAll("volumeBar") // 畫成交量bar chart
                .append("g")
                .data(volumeData)
                .enter().append("g")
                .attr("clip-path", "url(#clip)");

        //成交量bar顏色判斷
        chart.append("rect")
                .attr("class", "volumeBar")
                .attr("x", function (d) { return xScale(d.date); })
                .attr("height", function (d) {
                        return height - yVolume(d.volume);
                })
                .attr("y", function (d) {
                        return yVolume(d.volume);
                })
                .attr("width", xScale.bandwidth())
                .style("fill", function (d, i) { // 根據漲跌幅去決定成交量的顏色
                        if (data[i].change > 0) {
                                return "#FF0000"
                        } else if (data[i].change < 0) {
                                return "#00AA00"
                        } else {
                                return "#DDDDDD"
                        }
                });

        // 畫X軸 
        svg.selectAll("g.x.axis").call(xAxis.ticks(7).tickFormat(d3.timeFormat("%m/%d")).tickSize(-height, -height));

        //畫K線圖Y軸
        svg.selectAll("g.y.axis").call(yAxis.ticks(10).tickSize(-width, -width));

        //畫Ｋ線圖
        var state = svg.selectAll("g.candlestick")
                .attr("clip-path", "url(#candlestickClip)")
                .datum(data);
        state.call(candlestick)
                .each(function (d) {
                        dataArr = d;
                });
        //     console.log(data);
        //     console.log(testData);
        //     console.log(techan.indicator.sma().period(50)(data));//重要
        svg.select("g.sma.ma-0").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(5)(data)).call(sma0);
        svg.select("g.sma.ma-1").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(10)(data)).call(sma0);
        svg.select("g.ema.ma-2").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(20)(data)).call(sma0);
        svg.select("g.volume.axis").call(volumeAxis);
        svg.select("g.tradearrow").datum(trades).call(tradearrow);

        // 畫十字線並對他設定zoom function
        svg.append("g")
                .attr("class", "crosshair")
                .attr("width", width)
                .attr("height", height)
                .attr("pointer-events", "all")
                .call(crosshair)
                .call(zoom);

        //設定zoom的初始值
        zoomableInit = x.zoomable().clamp(false).copy();
        yInit = y.copy();
}

//設定當移動的時候要顯示的文字
function move(coords, index) {
        var i;
        for (i = 0; i < dataArr.length; i++) {
                if (coords.x === dataArr[i].date) {
                        svgText.text(d3.timeFormat("%Y/%m/%d")(coords.x) + "【開盤：" + dataArr[i].open + "】【高：" + dataArr[i].high + "】【低：" + dataArr[i].low + "】【收盤：" + dataArr[i].close + 
                        "】【漲跌：" + dataArr[i].change + "(" + dataArr[i].percentChange + "%)" + "】【成交量：" + dataArr[i].volume + 
                        "】【5MA：" + dataMaArr[i]["5MA"] + "】【10MA：" + dataMaArr[i]["10MA"] + "】【20MA：" + dataMaArr[i]["20MA"] + "】【60MA：" + dataMaArr[i]["60MA"]+"】"
                        ) 
                        ;
                }
        }
}

var rescaledX, rescaledY;
var t;
function zoomed() {

        //根據zoom去取得座標轉換的資料
        t = d3.event.transform;
        rescaledX = d3.event.transform.rescaleY(x);
        rescaledY = d3.event.transform.rescaleY(y);
        // y座標zoom
        yAxis.scale(rescaledY);
        candlestick.yScale(rescaledY);
        smac.yScale(rescaledY);
        sma0.yScale(rescaledY);
        sma1.yScale(rescaledY);
        ema2.yScale(rescaledY);
        tradearrow.yScale(rescaledY);
        // Emulates D3 behaviour, required for financetime due to secondary zoomable scale
        //K線圖 x zoom
        x.zoomable().domain(d3.event.transform.rescaleX(zoomableInit).domain());
        // 成交量 x  zoom
        xScale.range([0, width].map(d => d3.event.transform.applyX(d)));
        // y.range([0, height].map(d => d3.event.transform.applyY(d)));

        // 更新座標資料後，再重新畫圖
        redraw();
}



function redraw() {
        svg.select("g.candlestick").call(candlestick);
        svg.select("g.x.axis").call(xAxis);
        svg.select("g.y.axis").call(yAxis);
        svg.select("g.sma.ma-c.get").call(smac);
        svg.select("g.sma.ma-c.loss").call(smac);
        // svg.select("g.sma.ma-0").call(sma0);
        // svg.select("g.sma.ma-1").call(sma1);
        svg.select("g.ema.ma-2").call(ema2);
        svg.select("g.tradearrow").call(tradearrow);

        svg.selectAll("rect.volumeBar")
                .attr("x", function (d) { return xScale(d.date); })
                .attr("width", (xScale.bandwidth()));
}









