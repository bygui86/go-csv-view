package statics

const ViewTemplate = `
$(function () { setInterval({{ .ViewID }}_sync, {{ .Interval }}); });
function {{ .ViewID }}_sync() {
    $.ajax({
        type: "GET",
        url: "http://{{ .Address }}{{ .ViewPath }}",
        dataType: "json",
        success: function (result) {
            let opt = goecharts_{{ .ViewID }}.getOption();

            let x = opt.xAxis[0].data;
            x.push(result.time);
            if (x.length > {{ .MaxPoints }}) {
                x = x.slice(1);
            }
            opt.xAxis[0].data = x;

            for (let i = 0; i < result.values.length; i++) {
                let y = opt.series[i].data;
                y.push({ value: result.values[i] });
                if (y.length > {{ .MaxPoints }}) {
                    y = y.slice(1);
                }
                opt.series[i].data = y;

                goecharts_{{ .ViewID }}.setOption(opt);
            }
        }
    });
}`
