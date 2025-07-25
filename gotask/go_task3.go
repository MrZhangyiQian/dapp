type Student struct {
	Id    int `gorm:"primaryKey;autoIncrement;"`
	Name  string
	Age   int
	Grade string
}


func main() {
  db, err := gorm.Open("mysql", "root:123456@(localhost)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//
	// 自动迁移
	db.AutoMigrate(&dbtest.Student{})
	su1 := dbtest.Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	}
  // 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
	db.Create(&su1)
  
  //编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息
	var students []dbtest.Student
	db.Where("age > ?", 18).Find(&students)
	for _, students := range students {
		fmt.Printf("Name: %s, Age: %d, Grade: %s\n", students.Name, students.Age, students.Grade)
	}
  // 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	var suZhang = new(dbtest.Student)
	db.Where("name = ?", "张三").Find(&suZhang)
	fmt.Printf("Name: %s, Age: %d, Grade: %s", suZhang.Name, suZhang.Age, suZhang.Grade)
	suZhang.Grade = "四年级"
	db.Save(&suZhang)
  // 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	db.Where("age < ?", 15).Delete(&dbtest.Student{})
	var students []dbtest.Student
	db.Find(&students)
	for _, students := range students {
		fmt.Printf("Name: %s, Age: %d, Grade: %s\n", students.Name, students.Age, students.Grade)
	}
  }
// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务
CREATE TABLE accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '账户ID',
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00 COMMENT '账户余额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_balance (balance)
) ENGINE=InnoDB;

CREATE TABLE transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '交易ID',
    from_account_id BIGINT NOT NULL COMMENT '转出账户ID',
    to_account_id BIGINT NOT NULL COMMENT '转入账户ID',
    amount DECIMAL(15, 2) NOT NULL COMMENT '转账金额',
    transaction_time DATETIME NOT NULL COMMENT '交易时间',
    status ENUM('PENDING', 'SUCCESS', 'FAILED') DEFAULT 'SUCCESS' COMMENT '交易状态',
    
    FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    FOREIGN KEY (to_account_id) REFERENCES accounts(id),
    INDEX idx_time (transaction_time)
) ENGINE=InnoDB;

-- 1. 开始事务
START TRANSACTION;

-- 2. 声明变量存储账户余额
SET @balanceA = 0;
SET @balanceB = 0;

-- 3. 检查账户 A 的余额（并加锁防止并发操作）
SELECT balance INTO @balanceA 
FROM accounts 
WHERE id = A 
FOR UPDATE; -- 行级锁，阻止其他事务修改

-- 4. 检查账户 B 是否存在（也加锁）
SELECT balance INTO @balanceB 
FROM accounts 
WHERE id = B 
FOR UPDATE; 

-- 5. 验证账户余额和账户存在
IF @balanceA IS NULL THEN
    ROLLBACK;
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '转出账户不存在';
ELSEIF @balanceB IS NULL THEN
    ROLLBACK;
    SIGNAL SQLSTATE '45001' SET MESSAGE_TEXT = '转入账户不存在';
ELSEIF @balanceA < 100 THEN
    ROLLBACK;
    SIGNAL SQLSTATE '45002' SET MESSAGE_TEXT = '账户余额不足';
ELSE
    -- 6. 执行转账操作
    -- 6.1 从账户 A 扣除 100 元
    UPDATE accounts 
    SET balance = balance - 100 
    WHERE id = A;
    
    -- 6.2 向账户 B 增加 100 元
    UPDATE accounts 
    SET balance = balance + 100 
    WHERE id = B;
    
    -- 6.3 记录转账信息
    INSERT INTO transactions (from_account_id, to_account_id, amount, transaction_time)
    VALUES (A, B, 100, NOW());
    
    -- 7. 提交事务
    COMMIT;
    
    SELECT '转账成功' AS result;
END IF;

// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

// 自动迁移
	db.AutoMigrate(&dbtest.Employee{})
	emp1 := dbtest.Employee{
		Name: "张三",
		Department: "技术部",
		Salary: 5000,
	}
	emp2 := dbtest.Employee{
		Name: "李四",
		Department: "技术部",
		Salary: 6000,
	}
	db.Create(&emp1)
	db.Create(&emp2)
	var employees []dbtest.Employee
	db.Where("department = ?", "技术部").Find(&employees)
	//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	var employee dbtest.Employee
	db.Raw("SELECT * FROM employees ORDER BY salary DESC LIMIT 1").Scan(&employee)
	fmt.Printf("Name: %s, Department: %s, Salary: %d\n", employee.Name, employee.Department, employee.Salary)

//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
// 查询价格大于 50 的书籍
	var books []Book
	err = db.Select(&books, "SELECT * FROM books WHERE price > ?", 50)
	if err != nil {
		log.Fatalln("Query failed:", err)
	}

	// 打印结果
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
//假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//要求 ：
//使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//编写Go代码，使用Gorm创建这些模型对应的数据库表。
package ormtest

type Comment struct {
	Id      int    `gorm:"primary_key"`
	Content string `gorm:"type:text"`
	UserId  int
	PostId  int
	User    User
	Post    Post
}

func (c *Comment) TableName() string {
	return "comment"
}

package ormtest

type Post struct {
	Id       int    `gorm:"primary_key"`
	Title    string `gorm:"size:255"`
	Content  string
	UserId   int
	User     User
	Comments []Comment
}

func (p *Post) TableName() string {
	return "post"
}

package ormtest

type User struct {
	Id    int    `gorm:"primaryKey;autoIncrement;"`
	Name  string `gorm:"size:255"`
	Age   int
	Email string `gorm:"size:255;unique_index"`
	Posts []Post
}

func (User) TableName() string {
	return "user"
}


func main(){
db, err := gorm.Open("mysql", "root:123456@(localhost)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&ormtest.Comment{}, &ormtest.Post{}, &ormtest.User{})
}


// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func getUserPostsAndComments(db *gorm.DB, userID uint) ([]PostWithComments, error) {
	var results []PostWithComments

	rows, err := db.Table("post").
		Select("post.id as post_id, post.title as post_title, post.content as post_content, comment.id, comment.content, comment.user_id, comment.post_id").
		Joins("left join comment on post.id = comment.post_id").
		Where("post.user_id = ?", userID).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pwc PostWithComments
		var comment ormtest.Comment
		// Scan into the structs manually
		db.ScanRows(rows, &pwc)
		db.ScanRows(rows, &comment)

		// 判断是否已有该 Post，避免重复添加
		found := false
		for i, result := range results {
			if result.PostID == pwc.PostID {
				results[i].Comments = append(results[i].Comments, comment)
				found = true
				break
			}
		}

		if !found {
			pwc.Comments = []ormtest.Comment{comment}
			results = append(results, pwc)
		}
	}

	return results, nil
}

func main(){
db, err := gorm.Open("mysql", "root:123456@(localhost)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 查询用户 ID=1 的所有文章及评论
	userID := uint(1)
	posts, err := getUserPostsAndComments(db, userID)
	if err != nil {
		panic(err)
	}

	for _, post := range posts {
		fmt.Printf("Post: %s\n", post.PostTitle)
		for _, comment := range post.Comments {
			fmt.Printf("  Comment: %s (by user %d)\n", comment.Content, comment.UserId)
		}
	}
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段
package ormtest

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	Id       int    `gorm:"primary_key"`
	Title    string `gorm:"size:255"`
	Content  string
	UserId   int
	User     User
	Comments []Comment
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	// 查询用户
	var user User
	tx.Model(&User{}).Where("id = ?", p.UserId).First(&user)

	// 更新用户文章数
	user.PostCount++
	tx.Save(&user)

	return
}

func (p *Post) TableName() string {
	return "post"
}

func main() {
	db, err := gorm.Open("mysql", "root:123456@(localhost)/blog_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Post{})

	// 创建用户
	user := User{Name: "Alice"}
	db.Create(&user)

	// 创建文章
	post := Post{
		Title:   "My First Post",
		Content: "Hello World!",
		UserId:  user.ID,
	}
	db.Create(&post)

	// 查看用户文章数是否增加
	var updatedUser User
	db.First(&updatedUser, user.ID)
	fmt.Printf("User %s has %d posts\n", updatedUser.Name, updatedUser.PostCount)
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
package ormtest

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	Id      int    `gorm:"primary_key"`
	Content string `gorm:"type:text"`
	UserId  int
	PostId  int
	User    User
	Post    Post
}

func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 获取关联的文章
	var post Post
	tx.First(&post, c.PostId)

	// 查询该文章是否还有评论
	var count int
	tx.Model(&Comment{}).Where("post_id = ?", post.Id).Count(&count)

	// 如果没有评论了，更新文章的 HasComments 字段
	if count == 0 {
		post.HasComments = false
		tx.Save(&post)
	}

	return nil
}

func (c *Comment) TableName() string {
	return "comment"
}


func main() {
	db, err := gorm.Open("mysql", "root:123456@(localhost)/blog_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	// 创建用户
	user := User{Name: "Alice"}
	db.Create(&user)

	// 创建文章
	post := Post{
		Title:       "My First Post",
		Content:     "Hello World!",
		UserId:      user.ID,
		HasComments: true,
	}
	db.Create(&post)

	// 创建两个评论
	comment1 := Comment{
		Content: "Great article!",
		UserId:  user.ID,
		PostId:  post.Id,
	}
	db.Create(&comment1)

	comment2 := Comment{
		Content: "Nice job!",
		UserId:  user.ID,
		PostId:  post.Id,
	}
	db.Create(&comment2)

	// 删除第一个评论
	db.Delete(&comment1)

	// 再次删除第二个评论
	db.Delete(&comment2)

	// 查看文章的 HasComments 是否变为 false
	var updatedPost Post
	db.First(&updatedPost, post.Id)
	fmt.Printf("Post %s has comments: %v\n", updatedPost.Title, updatedPost.HasComments)
}
