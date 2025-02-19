function filterPosts() {
    let checked = document.querySelectorAll('input[name="categoryFilter"]:checked');
    let categoriesChecked = [];
    
    for (let i = 0; i < checked.length; i++) {
        categoriesChecked.push(checked[i].value);
    }

    let posts = document.querySelectorAll('article[id^="article_post_"]');
    for (let i = 0; i < posts.length; i++) {
        let post = posts[i];
        if (categoriesChecked.length == 0) {
            post.style.display = 'block';
            continue;
        }

        let postCategories = post.querySelector('.posts_categories').textContent.trim().split(' '); // Ensure proper splitting
        let matchCount = 0;

        // Check if post contains all selected categories
        for (let j = 0; j < categoriesChecked.length; j++) {
             
            if (postCategories.includes(categoriesChecked[j])) { 
                matchCount++;
            }
        }
        
        if (matchCount === categoriesChecked.length) {
            console.log(matchCount , categoriesChecked.length)
            post.style.display = 'block';
        } else {
            post.style.display = 'none';
        }
    }
}
