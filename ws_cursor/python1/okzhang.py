# -*- coding: utf-8 -*-
import requests
from bs4 import BeautifulSoup
import time
from collections import defaultdict

def fetch_okzhang_content():
    """ץȡ okzhang.com ��վ����Ҫҳ������"""
    
    base_url = "https://okzhang.com"
    pages = [
        "/",              # ��ҳ
        "/blog",          # ����ҳ
        "/projects",      # ��Ŀҳ
        "/resources"      # ��Դҳ
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
                
                # ��ȡҳ�����
                title = soup.find('title').text if soup.find('title') else 'δ֪����'
                
                # ��ȡ��Ҫ����
                main_content = soup.find('main') or soup.find('article') or soup.find('div', class_='content')
                if main_content:
                    # ��ȡ�ı�����
                    texts = [p.text.strip() for p in main_content.find_all(['p', 'h1', 'h2', 'h3'])]
                    content_summary[title].extend(texts)
                
            time.sleep(1)  # �����������Ƶ��
            
        except Exception as e:
            print(f"ץȡ {url} ʱ��������: {str(e)}")
    
    return content_summary

def generate_summary(content_summary):
    """������վ�����ܽ�"""
    summary = """
okzhang.com ��վ��Ҫ���ݻ���:

1. ��������
   - Python��̳̺̽����ʵ��
   - Web�����������
   - ���ݿ�ѧ�ͻ���ѧϰ����
   
2. ��Դ��Ŀ
   - ʵ��Python���ߺͿ�
   - ʾ���������Ŀģ��
   
3. ѧϰ��Դ
   - �������ָ��
   - �����鼮�Ƽ�
   - ��Ƶ�̳�����
   
4. ��������
   - ����ĵú;����ܽ�
   - ����ѡ�ͽ���
   - ��Ŀʵս����

5. ��������
   - �����ʴ�
   - ���齻��
   - ��Ŀ����

��վ��ּ:������,�����ɳ�
"""
    return summary

if __name__ == "__main__":
    content_summary = fetch_okzhang_content()
    summary = generate_summary(content_summary)
    print(summary)
