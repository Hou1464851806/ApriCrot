import sqlite3

def updateDB(exeLocation, electronVersion, electronArch):
    # 连接到SQLite数据库文件（如果不存在则会创建）
    conn = sqlite3.connect('electronDB.db')

    # 创建一个游标对象
    cursor = conn.cursor()

    # 创建应用信息表
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS electronApps (
            exeLocation TEXT PRIMARY KEY,
            electronVersion TEXT,
            electronArch TEXT
        );
    ''')

    # 插入应用信息到表中
    cursor.execute('''
        INSERT OR REPLACE INTO electronApps (exeLocation, electronVersion, electronArch)
        VALUES (?, ?, ?);
    ''', (exeLocation, electronVersion,  electronArch))

    # 提交更改
    conn.commit()

    # 关闭连接
    conn.close()

def checkDB():
    conn = sqlite3.connect('electronDB.db')
    cursor = conn.cursor()

    # 执行SQL查询
    query = "SELECT * FROM electronApps"
    cursor.execute(query)

    # 获取查询结果
    rows = cursor.fetchall()

    # 关闭连接
    conn.close()

    return rows
    