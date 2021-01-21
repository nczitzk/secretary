<p align="center">
    <a href="https://github.com/nczitzk/secretary" target="_blank">
        <img width="200" src="https://github.com/nczitzk/secretary/wiki/images/secretary.png" alt="secretary-logo">
    </a>
</p>

<p align="center">
    <a href="README.md">English</a> | 
    <a href="README-zh_CN.md">简体中文</a> | 正體中文
</p>

Secretary 是一個智慧排班助理，其能夠由定制的 XLSX（Microsoft Office Excel 2007 Document）範本快速生成值班表。

## 開始使用

建議下載包含完整的示例的 [壓縮包](https://github.com/nczitzk/secretary/releases)，並運行可執行檔，查看生成的 XLSX 檔。或者你也可以 pull 源碼，自行編譯。

如果一切順利，在 secretary 所在的目錄會列出兩個新的 XLSX 檔。一個是 secretary 生成的直接可用的值班表。另一個是所有職員可值班時間段的匯總表，清楚地顯示出每個指定的班次可以分配哪些職員。

![值班表](https://github.com/nczitzk/secretary/wiki/images/timetable.png)

![可用時間表](https://github.com/nczitzk/secretary/wiki/images/available-timetable.png)

## 工作流程

1. 製作一張在下麵步驟中會分發給職員的表格。建議在 Excel 中使用下拉式功能表，讓職員選擇是否接受該時間段的安排。

![空閒](https://github.com/nczitzk/secretary/wiki/images/all-free.gif)

2. 現在根據剛才創建的工作表中配置 JSON 檔（本示例中的 `settings.json`）中設置的 `__template_pattern` （例子中的 `{{Mon:1}}`）填入時間段識別字，以便 secretary 能正確理解職員的時間表。這也為下面的步驟中生成可用時間匯總表提供了範本。查看 [範本](https://github.com/nczitzk/secretary/wiki/Templates)。

![可用時間匯總表範本](https://github.com/nczitzk/secretary/wiki/images/mon-1.png)。

3. 派發時間表供職員填寫。

![職員填寫表格](https://github.com/nczitzk/secretary/wiki/images/free-occupied.gif)

4. 接下來需要製作值班表的範本。同樣，在需要 secretary 安排的儲存格中填寫時間段識別字。與可用時間匯總表不同，值班表的每個儲存格中只有一個職員的名字。這意味著你可以通過限制每個時間段的儲存格數量來確定每個班次的最多人數，甚至可以在值班表中為同一班次的職員指定不同的職位。是的，secretary 會儘量將職員分配到不同的崗位，而不是重複將同一職員分配到同一崗位。

![值班表範本](https://github.com/nczitzk/secretary/wiki/images/mon-2.png)

5. 在執行程式之前，請確認 JSON 檔（本示例中的 `settings.json`）中的配置已經正確設置。查看 [配置](https://github.com/nczitzk/secretary/wiki/Configurations)。

6. 運行 secretary。

![按兩下運行](https://github.com/nczitzk/secretary/wiki/images/double-click.gif)

## 更多

前往 [wiki](https://github.com/nczitzk/secretary/wiki) 瞭解更多關於 [參數](https://github.com/nczitzk/secretary/wiki/Parameters)、[範本](https://github.com/nczitzk/secretary/wiki/Templates) 和 [配置](https://github.com/nczitzk/secretary/wiki/Configurations) 的信息。

