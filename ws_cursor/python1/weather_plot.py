# -*- coding: utf-8 -*-
import pandas as pd
import matplotlib.pyplot as plt

# Data
data = {
    'month': ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12'],
    'huizhou': [14, 15, 18, 22, 26, 28, 29, 29, 28, 24, 20, 16],
    'zhengzhou': [1, 4, 10, 17, 23, 27, 28, 27, 22, 16, 9, 3],
    'harbin': [-18, -14, -4, 8, 16, 21, 23, 22, 15, 7, -5, -15]
}

# Create DataFrame
df = pd.DataFrame(data)

# Create line plot
plt.figure(figsize=(12, 6))
plt.plot(df['month'], df['huizhou'], marker='o', label='Huizhou')
plt.plot(df['month'], df['zhengzhou'], marker='s', label='Zhengzhou')
plt.plot(df['month'], df['harbin'], marker='^', label='Harbin')

# Set title and labels
plt.title('Temperature Trends of Three Cities', fontsize=14)
plt.xlabel('Month')
plt.ylabel('Temperature (C)')

# Add legend and grid
plt.legend()
plt.grid(True, linestyle='--', alpha=0.3)

# Save plot
plt.savefig('weather_trend.png', dpi=300, bbox_inches='tight')
plt.close() 