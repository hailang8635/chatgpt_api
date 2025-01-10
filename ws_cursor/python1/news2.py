# -*- coding: utf-8 -*-
import requests
from bs4 import BeautifulSoup
import pandas as pd
from datetime import datetime

def get_baidu_hot():
    # 获取百度热搜页面
    url = "https://top.baidu.com/board?tab=realtime"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    
    try:
        response = requests.get(url, headers=headers)
        response.encoding = 'utf-8'  # 明确指定响应的编码
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # 存储热搜数据的列表
        hot_items = []
        
        # 获取热搜条目
        items = soup.find_all('div', class_='category-wrap_iQLoo')
        
        for item in items:
            # 获取排名
            rank = len(hot_items) + 1
            
            # 获取标题
            title = item.find('div', class_='c-single-text-ellipsis').text.strip()
            
            # 获取热度值
            hot_value = item.find('div', class_='hot-index_1Bl1a').text.strip()
            
            hot_items.append({
                'rank': rank,
                'title': title,
                'hot_value': hot_value,
                'datetime': datetime.now().strftime('%Y-%m-%d %H:%M:%S')
            })
        
        # 转换为DataFrame
        df = pd.DataFrame(hot_items)
        
        # 保存到CSV文件，如果文件不存在则创建，存在则追加
        df.to_csv('baidu_hot.csv', mode='a', header=not pd.io.common.file_exists('baidu_hot.csv'), index=False, encoding='utf-8-sig')
        
        print(f"成功获取{len(hot_items)}条百度热搜数据")
        return df
        
    except Exception as e:
        print(f"获取百度热搜失败: {str(e)}")
        return None

if __name__ == "__main__":
    get_baidu_hot()
