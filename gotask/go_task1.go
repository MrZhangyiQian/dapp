// 136. 只出现一次的数字：
func singleNumber(nums []int) int {
    // 1 2 2 3 1
    var freq map[int]bool
    freq = make(map[int]bool)
    for _,value := range nums{
        if(freq[value]){
            delete(freq,value)
        }else{
            freq[value] = true
        }
    }

    for key :=  range freq{
        return key
    } 
    return -1
}

func singleNumber(nums []int) int {
    // 1 2 2 3 1   满足交换率结合律， 
    index := 0
    for _;num := range(nums){
      index ^= num
}
 return index
}

// 有效的括号 
func isValid(s string) bool {
    // 检查长度是否为奇数
    if len(s)%2 == 1 {
        return false
    }
    
    // 括号映射：右括号 -> 左括号
    pairs := map[byte]byte{
        ')': '(',
        ']': '[',
        '}': '{',
    }
    
    // 使用栈保存左括号
    stack := []byte{}
    
    for i := 0; i < len(s); i++ {
        // 如果当前字符是右括号（pairs[s[i]]存在）
        if char, exists := pairs[s[i]]; exists {
            // 检查栈是否为空或栈顶不匹配
            if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
                return false
            }
            // 弹出匹配的左括号
            stack = stack[:len(stack)-1]
        } else {
            // 否则（左括号或其他字符）压入栈
            stack = append(stack, s[i])
        }
    }
    
    // 检查所有括号是否都已匹配
    return len(stack) == 0
}

func isValid(s string) bool {
	stack := make([]byte, 1000000)
	top := -1
	for _, c := range s {
		if c == '[' || c == '{' || c == '(' {
			top++
			stack[top] = byte(c)
		} else if top > -1 && match(stack[top], byte(c)) {
			top--
		} else {
			return false
		}
	}
	return top == -1
}

func match(a, b byte) bool {
	return (a == '[' && b == ']') || (a == '{' && b == '}') || (a == '(' && b == ')')
}


// 最长公共前缀
func longestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }
    prefix := strs[0]

    for i:=1; i<len(strs); i++ {
        // 字符串和前缀的公共长度
        j := 0
        for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
            j++
        }
        // 更新前缀
        prefix = prefix[:j]
        if prefix == "" {
            return prefix
        }
    }
    return prefix
}

// 26. 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
    // 1 1 2 2 3 4 5
    // 返回的是新数组长度,原地删除重复元素
    if len(nums) == 0 {
         return 0
    }
    fast := 1
    slow := 1
    for fast < len(nums) {
        if nums[fast] != nums[fast-1] {
            nums[slow] = nums[fast]
            slow++
        }
        fast++
    }
    return slow
    
}

// 加一
func plusOne(digits []int) []int {
        for i := len(digits) - 1; i >= 0; i-- {
        digits[i]++
        digits[i] = digits[i] % 10
        if (digits[i] != 0){
           return digits
        }
    }
    digits = make([]int,len(digits)+1)
    digits[0] = 1
    return digits
}

// 两数之和
func twoSum(nums []int, target int) []int {
    hp := map[int]int{}
    for i, x:= range nums{
        if p,ok := hp[target-x]; ok{
            return []int{p,i}
        }
        hp[x] = i
    }
    return nil
    
}