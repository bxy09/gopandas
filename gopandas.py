import pandas as pd
import numpy as np
import gopandas_pb2 as gpb
import pytz
import datetime


class GoPandas:
    def __init__(self):
        return

    @staticmethod
    def panel_from_gopanel(go_panel):
        tz = pytz.timezone("Asia/Shanghai")
        dates = [tz.localize(datetime.datetime.fromtimestamp(unix/1000000000)) for unix in go_panel.dates]
        array = np.reshape(newshape=(len(go_panel.dates), len(go_panel.secondary), len(go_panel.thirdly)),
                           a=go_panel.data, order='C')
        return pd.Panel(data=array, items=dates, major_axis=go_panel.secondary, minor_axis=go_panel.thirdly)

    # parse pandas panel (index: unix_time,string,string) with proto_buf
    @staticmethod
    def panel_from_protobuf(string):
        go_panel = gpb.TimePanel()
        go_panel.ParseFromString(string)
        return GoPandas.panel_from_gopanel(go_panel)

    # serialize pandas panel (index: unix_time,string,string) with proto_buf
    @staticmethod
    def panel_to_protobuf(panel):
        go_panel = gpb.TimePanel()
        axes = panel.axes
        dates = axes[0].tolist()
        unix = [int(date.strftime("%s"))*1000000000 for date in dates]
        go_panel.dates.extend(unix)
        go_panel.secondary.extend(axes[1].tolist())
        go_panel.thirdly.extend(axes[2].tolist())
        go_panel.data.extend(panel.values.reshape((len(go_panel.dates)*len(go_panel.secondary)*len(go_panel.thirdly)),
                                                  order='C').tolist())
        return go_panel.SerializeToString()

if __name__ == "__main__":
    def test():
        shtz = pytz.timezone("Asia/Shanghai")
        start = datetime.datetime(2007,1,1)
        start = shtz.localize(start)
        end = datetime.datetime(2008,1,1)
        end = shtz.localize(end)

        p = pd.Panel(
            data=[[[0,1],[2,3]],[[4,5],[6,7]]],
            items=[start,end],
            major_axis=['000001.SZ','000002.SZ'],
            minor_axis=['CLOSE','OPEN'],
        )
        print p
        print p.values
        string = GoPandas.panel_to_protobuf(p)
        pp = GoPandas.panel_from_protobuf(string)

        print pp
        print pp.values

    test()
