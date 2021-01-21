<p align="center">
    <a href="https://github.com/nczitzk/secretary" target="_blank">
        <img width="200" src="https://github.com/nczitzk/secretary/wiki/images/secretary.png" alt="secretary-logo">
    </a>
</p>

<p align="center">
    <a href="README.md">English</a> | 简体中文 | 
    <a href="README-zh_TW.md">正體中文</a>
</p>

Secretary 是一个智能排班助理，其能够由定制的 XLSX（Microsoft Office Excel 2007 Document）模板快速生成值班表。

## 开始使用

建议下载包含完整的示例的 [压缩包](https://github.com/nczitzk/secretary/releases)，并运行可执行文件，查看生成的 XLSX 文件。或者你也可以 pull 源码，自行编译。

如果一切顺利，在 secretary 所在的目录会列出两个新的 XLSX 文件。一个是 secretary 生成的直接可用的值班表。另一个是所有职员可值班时间段的汇总表，清楚地显示出每个指定的班次可以分配哪些职员。

![值班表](https://github.com/nczitzk/secretary/wiki/images/timetable.png)

![可用时间表](https://github.com/nczitzk/secretary/wiki/images/available-timetable.png)

## 工作流程

1. 制作一张在下面步骤中会分发给职员的表格。建议在 Excel 中使用下拉菜单，让职员选择是否接受该时间段的安排。

![空闲](https://github.com/nczitzk/secretary/wiki/images/all-free.gif)

2. 现在根据刚才创建的工作表中配置 JSON 文件（本示例中的 `settings.json`）中设置的 `__template_pattern` （例子中的 `{{Mon:1}}`）填入时间段标识符，以便 secretary 能正确理解职员的时间表。这也为下面的步骤中生成可用时间汇总表提供了模板。查看 [模板](https://github.com/nczitzk/secretary/wiki/Templates)。

![可用时间汇总表模板](https://github.com/nczitzk/secretary/wiki/images/mon-1.png)。

3. 派发时间表供职员填写。

![职员填写表格](https://github.com/nczitzk/secretary/wiki/images/free-occupied.gif)

4. 接下来需要制作值班表的模板。同样，在需要 secretary 安排的单元格中填写时间段标识符。与可用时间汇总表不同，值班表的每个单元格中只有一个职员的名字。这意味着你可以通过限制每个时间段的单元格数量来确定每个班次的最多人数，甚至可以在值班表中为同一班次的职员指定不同的职位。是的，secretary 会尽量将职员分配到不同的岗位，而不是重复将同一职员分配到同一岗位。

![值班表模板](https://github.com/nczitzk/secretary/wiki/images/mon-2.png)

5. 在执行程序之前，请确认 JSON 文件（本示例中的 `settings.json`）中的配置已经正确设置。查看 [配置](https://github.com/nczitzk/secretary/wiki/Configurations)。

6. 运行 secretary。

![双击运行](https://github.com/nczitzk/secretary/wiki/images/double-click.gif)

## 更多

前往 [wiki](https://github.com/nczitzk/secretary/wiki) 了解更多关于 [参数](https://github.com/nczitzk/secretary/wiki/Parameters)、[模板](https://github.com/nczitzk/secretary/wiki/Templates) 和 [配置](https://github.com/nczitzk/secretary/wiki/Configurations) 的信息。