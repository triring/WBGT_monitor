python plotGraph.py WBGT_data.csv Temperature_Humidity.png --metrics Temperature Humidity --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 10 80 --grid
python plotGraph.py WBGT_data.csv Pressure.png --metrics Pressure --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --grid
python plotGraph.py WBGT_data.csv HeatIndex.png --metrics HeatIndex --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 0 35 --grid
python plotGraph.py WBGT_data.csv DiscomfortIndex.png --metrics DiscomfortIndex --xlim_start "2025-08-04 08:30" --xlim_end "2025-08-04 17:00" --ylim 50 100 --grid
