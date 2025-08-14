import pandas as pd
import matplotlib.pyplot as plt
import argparse
from datetime import datetime
from matplotlib.dates import DateFormatter
import matplotlib.dates as mdates

# 実行サンプル
# python ms_plot.py WBGT20250804.csv Temperature_Humidity.png --metrics Temperature Humidity --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 10 80 --grid
# python ms_plot.py WBGT20250804.csv Pressure.png --metrics Pressure --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --grid
# python ms_plot.py WBGT20250804.csv HeatIndex.png --metrics HeatIndex --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 0 35 --grid
# python ms_plot.py WBGT20250804.csv DiscomfortIndex.png --metrics DiscomfortIndex --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 50 80 --grid


# コマンドライン引数の設定
parser = argparse.ArgumentParser(description='CSVデータからグラフを作成してPNGで保存')
parser.add_argument('csv_file', help='読み込むCSVファイル名')
parser.add_argument('output_file', help='保存するPNGファイル名')
parser.add_argument('--xlim_start', help='横軸の開始時刻 (YYYY-MM-DD HH:MM)', default=None)
parser.add_argument('--xlim_end', help='横軸の終了時刻 (YYYY-MM-DD HH:MM)', default=None)
parser.add_argument('--ylim', nargs=2, type=float, help='縦軸の範囲 (min max)', default=None)
parser.add_argument('--grid', action='store_true', help='グリッドを表示する')
parser.add_argument('--metrics', nargs='+', choices=['Temperature', 'Humidity', 'Pressure', 'HeatIndex', 'DiscomfortIndex'],
                    help='表示する項目を指定（複数指定可）', default=['Temperature', 'Humidity', 'Pressure', 'HeatIndex', 'DiscomfortIndex'])

args = parser.parse_args()

# CSV読み込み
df = pd.read_csv(args.csv_file, header=None,
                 names=['Date', 'Time', 'Temperature', 'Humidity', 'Pressure', 'HeatIndex', 'DiscomfortIndex'])

# 日付と時刻を結合してdatetime型に変換
df['Datetime'] = pd.to_datetime(df['Date'] + ' ' + df['Time'])

# 横軸の範囲指定
if args.xlim_start:
    start_time = datetime.strptime(args.xlim_start, '%Y-%m-%d %H:%M')
    df = df[df['Datetime'] >= start_time]
if args.xlim_end:
    end_time = datetime.strptime(args.xlim_end, '%Y-%m-%d %H:%M')
    df = df[df['Datetime'] <= end_time]

# 測定日付の取得（最初の行の日付）
measurement_date = df['Date'].iloc[0]

# グラフ描画
plt.figure(figsize=(12, 6))
for metric in args.metrics:
    plt.plot(df['Datetime'], df[metric], label=metric)

plt.xlabel('Time')
plt.ylabel('Values')
plt.title(f'Measurements on {measurement_date}')
plt.legend()
if args.grid:
    plt.grid(True)
if args.ylim:
    plt.ylim(args.ylim)

# X軸のフォーマットを時刻のみに変更
plt.gca().xaxis.set_major_formatter(DateFormatter('%H:%M'))
plt.xticks(rotation=45)
plt.tight_layout()
plt.savefig(args.output_file)
