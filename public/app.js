//讀取資料列表-動態產生按鈕
$('#tax,#to_TWD').on('change',function () {
        
        toRealPrice();
})


var file_index = "data/index.csv";
var list_index;
d3.csv(file_index, function (error, data) {
        list_index = data;
        data.forEach(function (v,i) {
                $("#list").append(`<button onclick="loadJSON('${v.filename}','ma.csv','${v.action}')">${v.filename}</button>`);
        });
});


var margin = { top: 20, right: 50, bottom: 30, left: 60 },
        width = 1300 - margin.left - margin.right,
        height = 700 - margin.top - margin.bottom;

// 設定時間格式
var parseDate = d3.timeParse("%Y/%m/%d");
var parseTime = d3.timeParse("%H:%M:%S");
var parseDateTime = d3.timeParse("%Y/%m/%d %H:%M:%S");
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
        .format(d3.timeFormat('%Y-%m-%d %H:%M'))
        .translate([0, height]);

// 設定十字線
var crosshair = techan.plot.crosshair()
        .xScale(x)
        .yScale(crosshairY)
        .xAnnotation(timeAnnotation)
        .yAnnotation(ohlcAnnotation)
        .on("move", move);

// 設定文字區域
var textSvg = d3.select("#graph").append("svg")
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
var svg = d3.select("#graph")
        .append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .attr("pointer-events", "all")
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");


// var dataArr = [];
var dataMaArr_o = [];
var dataBuySellArr_o = [];
var g_data = [];//日期篩選後的k棒
var g_volumeData = [];//日期篩選後的量
var path = "data/";


// loadJSON("08:45_60min.csv", "ma.csv", "buyAndSell.csv", "date");


function loadJSON(k_file, ma_file, action_file) {
        var date_start = $('#start_year').val() + $('#start_month').val();
        var date_end = $('#end_year').val() + $('#end_month').val();
        date_start = monthDate(date_start)
        date_end = monthDate(date_end)
        date_end.setMonth(date_end.getMonth() + 1)
        svg.selectAll("*").remove(); // 切換不同資料需要重新畫圖，因此需要先清除原先的圖案
        d3.csv(ma_file, function (error, data) {
                dataMaArr_o = data;//匯入ma資料
        });

        d3.csv(path + action_file, function (error, data) {
                dataBuySellArr_o = data;//匯入買賣資料
        });
        d3.csv(path + k_file, function (error, data) {
                var accessor = candlestick.accessor();
                var jsonData = data;

                data = jsonData.map(function (d) { // 設定data的格式
                        return {
                                date: parseDateTime(d["Date"] + ' ' + d["Time"]),
                                open: +d["Open"],
                                high: +d["High"],
                                low: +d["Low"],
                                close: +d["Close"],
                                volume: +d["TotalVolume"]
                        };
                }).sort(function (a, b) { return d3.ascending(accessor.d(a), accessor.d(b)); });


                var volumeData = jsonData.map(function (d) {
                        return {
                                date: parseDateTime(d["Date"] + ' ' + d["Time"]),
                                volume: d["TotalVolume"]
                        }
                });
                //     .reverse();



                svg.append("g")
                        .attr("class", "candlestick");
                svg.append("g")
                        .attr("class", "sma ma-c get");
                svg.append("g")
                        .attr("class", "sma ma-c loss");
                svg.append("g")
                        .attr("class", "sma ma-0");
                svg.append("g")
                        .attr("class", "sma ma-1");
                svg.append("g")
                        .attr("class", "ema ma-2");
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
                g_data = data.filter(function (value) {//k棒
                        return date_start <= value.date && value.date <= date_end;
                });
                g_volumeData = volumeData.filter(function (value) {//成交量
                        return date_start <= value.date && value.date <= date_end;
                });

                dataMaArr = dataMaArr_o.filter(function (value) {//MA數字資料
                        return date_start <= parseDate(value['交易日期']) && parseDate(value['交易日期']) <= date_end;
                });
                dataBuySellArr = dataBuySellArr_o.filter(function (value) {//
                        return date_start <= parseDateTime(value['成交時間']) && parseDateTime(value['成交時間']) <= date_end;
                });
                // draw(data.slice(start, end), g_volumeData.slice(start, end));
                draw();
                drawTable();
                drawTable2();
                toRealPrice();
        });
}
function drawTable2() {
        date_start = $('#start_year').val() + $('#start_month').val();
        date_start = monthDate(date_start);

        
        data_month = [];
        data_year = [];
        for (let i = 0; i < 12; i++) {
                row = {"年份":"xxxx","year_count":0,"期間":(i+1)+'月',"平均獲利":0,"毛利":0,"毛損":0,"獲利因子":0,"平均交易數量":0,"勝率":"0%","win":0,"loss":0,"count":0,"total_price":0};
                data_month[i] = row;
        }

        dataBuySellArr_o.forEach(function (v,i) {
                let row = {"年份":"xxxx","year_count":0,"期間":'月',"平均獲利":0,"毛利":0,"毛損":0,"獲利因子":0,"平均交易數量":0,"勝率":"0%","win":0,"loss":0,"count":0,"total_price":0};

                if(v['類型'] == "sell"){//
                        month_key = parseDateTime(v['成交時間']).getMonth();
                        year = parseDateTime(v['成交時間']).getFullYear();
                        if(data_month[month_key]["年份"] != year){
                                data_month[month_key]["年份"] = year;
                                data_month[month_key]["year_count"] += 1;
                        }
                        data_month[month_key]["total_price"] += parseInt(v['獲利']);
                        data_month[month_key]["count"] += 1;
                        if(v['獲利'] > 0){
                                data_month[month_key]["win"] += 1;
                                data_month[month_key]["毛利"] += parseInt(v['獲利']);
                        }
                        if(v['獲利'] < 0){
                                data_month[month_key]["loss"] += 1;
                                data_month[month_key]["毛損"] += parseInt(v['獲利']);
                        }
                }
                if(v['類型'] == "sell"){//年份資料
                        month = parseDateTime(v['成交時間']).getMonth();
                        year_key = parseDateTime(v['成交時間']).getFullYear();
                        if(data_year[year_key] === undefined){
                                data_year[year_key] = row;
                        }
                        if(data_year[year_key]["年份"] != year_key){
                                data_year[year_key]["年份"] = year_key;
                        }
                        // data_year[year_key]["year_count"] += 1;
                        data_year[year_key]["total_price"] += parseInt(v['獲利']);
                        data_year[year_key]["count"] += 1;
                        if(v['獲利'] > 0){
                                data_year[year_key]["win"] += 1;
                                data_year[year_key]["毛利"] += parseInt(v['獲利']);
                        }
                        if(v['獲利'] < 0){
                                data_year[year_key]["loss"] += 1;
                                data_year[year_key]["毛損"] += parseInt(v['獲利']);
                        }
                }
        
        });
        table = `<table>`;
        // 期間,平均獲利,毛利,毛損,獲利因子,平均交易數量,勝率
        table += `<tr>`;
        table += `<th>期間</th>`;
        table += `<th>平均獲利</th>`;
        table += `<th>毛利</th>`;
        table += `<th>毛損</th>`;
        table += `<th>獲利因子</th>`;
        table += `<th>平均交易數量</th>`;
        table += `<th>贏</th>`;
        table += `<th>輸</th>`;
        table += `<th>勝率</th>`;
        table += `</tr>`;
        for (let i = 0; i < 12; i++) {
                let price = (data_month[i]['total_price']/data_month[i]['year_count']).toFixed(2);
                let win_price = (data_month[i]['毛利']/data_month[i]['year_count']).toFixed(2);
                let loss_price = (data_month[i]['毛損']/data_month[i]['year_count']).toFixed(2);
                table += `<tr>`;
                table += `<td>${data_month[i]['期間']}</td>`;
                table += `<td class="price ${price > 0 ? 'price_green':'price_red'}" data-price="${price}">${price}</td>`;
                table += `<td class="price ${win_price > 0 ? 'price_green':'price_red'}" data-price="${win_price}">${win_price}</td>`;
                table += `<td class="price ${loss_price > 0 ? 'price_green':'price_red'}" data-price="${loss_price}">${loss_price}</td>`;
                table += `<td>${(data_month[i]['毛利']/data_month[i]['毛損'] * -1).toFixed(2)}</td>`;
                table += `<td>${(data_month[i]['count']/data_month[i]['year_count']).toFixed(2)}</td>`;
                table += `<td>${(data_month[i]['win']/data_month[i]['year_count']).toFixed(2)}</td>`;
                table += `<td>${(data_month[i]['loss']/data_month[i]['year_count']).toFixed(2)}</td>`;
                table += `<td>${(data_month[i]['win'] / data_month[i]['count'] *100).toFixed(2)} %</td>`;
        
                table += `</tr>`;
        }
        table += `</table>`;
        $('#table_month').html(table);
        table = `<table>`;


        // 期間,平均獲利,毛利,毛損,獲利因子,平均交易數量,勝率
        table += `<tr>`;
        table += `<th>期間</th>`;
        table += `<th>獲利</th>`;
        table += `<th>毛利</th>`;
        table += `<th>毛損</th>`;
        table += `<th>獲利因子</th>`;
        table += `<th>交易數量</th>`;
        table += `<th>贏</th>`;
        table += `<th>輸</th>`;
        table += `<th>勝率</th>`;
        table += `</tr>`;

        data_year.forEach(function (v,i) {
                let price = (v['total_price']).toFixed(2);
                let win_price = (v['毛利']).toFixed(2);
                let loss_price = (v['毛損']).toFixed(2);
                table += `<tr>`;
                table += `<td>${v['年份']}</td>`;
                table += `<td class="price ${price > 0 ? 'price_green':'price_red'}" data-price="${price}">${price}</td>`;
                table += `<td class="price ${win_price > 0 ? 'price_green':'price_red'}" data-price="${win_price}">${win_price}</td>`;
                table += `<td class="price ${loss_price > 0 ? 'price_green':'price_red'}" data-price="${loss_price}">${loss_price}</td>`;
                table += `<td>${(v['毛利']/v['毛損'] * -1).toFixed(2)}</td>`;
                table += `<td>${(v['count']).toFixed(2)}</td>`;
                table += `<td>${(v['win']).toFixed(2)}</td>`;
                table += `<td>${(v['loss']).toFixed(2)}</td>`;
                table += `<td>${(v['win'] / v['count'] *100).toFixed(2)} %</td>`;
        
                table += `</tr>`;
        });
        table += `</table>`;
        $('#table_year').html(table);
        table = `<table>`;
}




function toRealPrice() {
        tax = $('#tax').val();//滑價
        to_TWD = $('#to_TWD').val();//一點價格
        $('.price').each(function (i,v) {
                price = $(v).data('price');
                lots = $(v).data('lots') != undefined ?  $(v).data('lots'): 1;
                if (price != 0){
                        price = price*to_TWD - 2*tax*lots;//每次交易算一次(TODO還沒計算口數)
                        $(v).text(price.toFixed(0));
                }
        })
}

function drawTable() {
        
        table = `<table>`;
        // 交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)
        table += `<tr>`;
        table += `<th>交易編號</th>`;
        table += `<th>委託單編號</th>`;
        table += `<th>類型</th>`;
        table += `<th>訊號</th>`;
        table += `<th style="width:210px">成交時間</th>`;
        table += `<th>價格</th>`;
        table += `<th>數量</th>`;
        table += `<th>獲利</th>`;
        table += `<th>累積獲利</th>`;
        table += `<th>最大可能獲利</th>`;
        table += `<th>最大可能虧損</th>`;
        table += `</tr>`;
        dataBuySellArr.forEach(function (v,i) {
                table += `<tr>`;
                table += `<td>${v['交易編號']}</td>`;
                table += `<td>${v['委託單編號']}</td>`;
                table += `<td>${v['類型']}</td>`;
                table += `<td>${v['訊號']}</td>`;
                table += `<td>${v['成交時間']}</td>`;
                table += `<td>${v['價格']}</td>`;
                table += `<td>${v['數量']}</td>`;
                table += `<td class="price ${v['獲利'] != 0 ? v['獲利'] >0 ? 'price_green' : 'price_red' : ''}" data-price="${v['獲利']}" data-lots="${v['數量']}">${v['獲利']}</td>`;
                table += `<td class="price ${v['累積獲利'] != 0 ? v['累積獲利'] >0 ? 'price_green' : 'price_red' : ''}" data-price="${v['累積獲利']}" data-lots="${v['數量']}">${v['累積獲利']}</td>`;
                table += `<td class="price ${v['最大可能獲利'] != 0 ? v['最大可能獲利'] >0 ? 'price_green' : 'price_red' : ''}" data-price="${v['最大可能獲利']}" data-lots="${v['數量']}">${v['最大可能獲利']}</td>`;
                table += `<td class="price ${v['最大可能虧損'] != 0 ? v['最大可能虧損'] >0 ? 'price_green' : 'price_red' : ''}" data-price="${v['最大可能虧損']}" data-lots="${v['數量']}">${v['最大可能虧損']}</td>`;
                table += `</tr>`;
        });

        table += `</table>`;
        $('#table_row').html(table);
        table = '';
        
}

function draw() {
        // 設定domain，決定各座標所用到的資料
        x.domain(g_data.map(candlestick.accessor().d));
        y.domain(techan.scale.plot.ohlc(g_data, candlestick.accessor()).domain());
        xScale.domain(g_volumeData.map(function (d) { return d.date; }))
        yVolume.domain(techan.scale.plot.volume(g_data).domain());


        var trades = dataBuySellArr.map(function (d) {
                return {
                        date: parseDateTime(d["成交時間"]),
                        type: d["類型"].trim(),
                        price: parseInt(d["價格"]),
                }
        });
        var get_line = [];
        var loss_line = [];
        var temp = [];
        dataBuySellArr.forEach(function (d) {
                temp.push(d)
                if (d["獲利"] > 0) {
                        temp.push({ date: null, value: 0 });
                        get_line = get_line.concat(temp);
                        temp = [];
                }
                if (d["獲利"] < 0) {
                        temp.push({ date: null, value: 0 });
                        loss_line = loss_line.concat(temp);
                        temp = []
                }
        });
        var get_line = get_line.map(function (d) {
                return {
                        date: parseDateTime(d["成交時間"]) ?? null,
                        value: d["價格"] === undefined ? null : parseInt(d["價格"])
                }
        });
        var loss_line = loss_line.map(function (d) {
                return {
                        date: parseDateTime(d["成交時間"]) ?? null,
                        value: d["價格"] === undefined ? null : parseInt(d["價格"]),
                }
        });
        // console.log(get_line);
        // console.log(loss_line);
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
                .data(g_volumeData)
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
                        if (g_data[i].change > 0) {
                                return "#FF0000"
                        } else if (g_data[i].change < 0) {
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
                .datum(g_data);

        state.call(candlestick);
        //         .each(function (d) {
        //                 dataArr = d;
        //         });


  
        svg.select("g.sma.ma-0").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(5)(g_data)).call(sma0);
        svg.select("g.sma.ma-1").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(10)(g_data)).call(sma0);
        svg.select("g.ema.ma-2").attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(20)(g_data)).call(sma0);
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
        for (i = 0; i < g_data.length; i++) {
                if (coords.x === g_data[i].date) {
                        svgText.text(d3.timeFormat("%Y/%m/%d %H:%M")(coords.x) + "【開盤：" + g_data[i].open + "】【高：" + g_data[i].high + "】【低：" + g_data[i].low + "】【收盤：" + g_data[i].close +
                                "】【成交量：" + g_data[i].volume + "】"
                        );
                }
        }
}

var rescaledX, rescaledY;
var t;
function zoomed() {

        //根據zoom去取得座標轉換的資料
        // t = d3.event.transform;
        rescaledX = d3.event.transform.rescaleY(x);
        rescaledY = d3.event.transform.rescaleY(y);
        // y座標zoom

        // yAxis.scale(rescaledY);
        // candlestick.yScale(rescaledY);
        // smac.yScale(rescaledY);
        // sma0.yScale(rescaledY);
        // sma1.yScale(rescaledY);
        // ema2.yScale(rescaledY);
        // tradearrow.yScale(rescaledY);

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
        svg.select("g.sma.ma-0").call(sma0);
        svg.select("g.sma.ma-1").call(sma1);
        svg.select("g.ema.ma-2").call(ema2);
        svg.select("g.tradearrow").call(tradearrow);

        svg.selectAll("rect.volumeBar")
                .attr("x", function (d) { return xScale(d.date); })
                .attr("width", (xScale.bandwidth()));
}

function MA(k) {
        if(k == 0){//全關
                svg.select("g.sma.ma-0").remove();
                svg.select("g.sma.ma-1").remove();
                svg.select("g.ema.ma-2").remove();
                return;
        }
        if (k == 5) {
                MAname = "g.sma.ma-0";
                MAname2 = "sma ma-0";
        }
        if (k == 10) {
                MAname = "g.sma.ma-1";
                MAname2 = "sma ma-1";
        }
        if (k == 20) {
                MAname = "g.ema.ma-2";
                MAname2 = "ema ma-2";
        }
        if (svg.select(MAname).size()) {
                svg.select(MAname).remove();
        } else {
                svg.append("g").attr("class", MAname2);
                svg.select(MAname).attr("clip-path", "url(#candlestickClip)").datum(techan.indicator.sma().period(k)(g_data)).call(sma0);
        }
}








