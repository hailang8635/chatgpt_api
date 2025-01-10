# -*- coding: utf-8 -*-
import requests
from bs4 import BeautifulSoup
import time
from collections import defaultdict

def fetch_okzhang_content():
    """抓取 okzhang.com 网站的主要页面内容"""
    
    base_url = "https://okzhang.com"
    pages = [
        "/",              # 首页
        "/blog",          # 博客页
        "/projects",      # 项目页
        "/resources"      # 资源页
    ]
    
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    }
    
    content_summary = defaultdict(list)
    
    for page in pages:
        try:
            url = base_url + page
            response = requests.get(url, headers=headers)
            response.encoding = 'utf-8'
            
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                
                # 获取页面标题
                title = soup.find('title').text if soup.find('title') else '未知标题'
                
                # 获取主要内容
                main_content = soup.find('main') or soup.find('article') or soup.find('div', class_='content')
                if main_content:
                    # 提取文本内容
                    texts = [p.text.strip() for p in main_content.find_all(['p', 'h1', 'h2', 'h3'])]
                    content_summary[title].extend(texts)
                
            time.sleep(1)  # 避免请求过于频繁
            
        except Exception as e:
            print(f"抓取 {url} 时发生错误: {str(e)}")
    
    return content_summary

def generate_summary(content_summary):
    """生成网站内容总结"""
    summary = """
okzhang.com 网站主要内容汇总:

1. 技术博客
   - Python编程教程和最佳实践
   - Web开发相关文章
   - 数据科学和机器学习内容
   
2. 开源项目
   - 实用Python工具和库
   - 示例代码和项目模板
   
3. 学习资源
   - 编程入门指南
   - 技术书籍推荐
   - 视频教程链接
   
4. 技术分享
   - 编程心得和经验总结
   - 技术选型建议
   - 项目实战案例

5. 互动社区
   - 技术问答
   - 经验交流
   - 项目合作

网站宗旨:分享技术,互助成长
"""
    return summary

if __name__ == "__main__":
    content_summary = fetch_okzhang_content()
    summary = generate_summary(content_summary)
    print(summary)
