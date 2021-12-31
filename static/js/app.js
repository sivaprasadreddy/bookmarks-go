new Vue({
    el: '#app',
    delimiters: ['${', '}'],
    data: {
        bookmarks: [],
        newBookmark: {}
    },
    created: function () {
        this.loadBookmarks();
    },
    methods: {
        loadBookmarks() {
            let self = this;
            $.getJSON("/api/bookmarks", function (data) {
                self.bookmarks = data
            });
        },

        saveBookmark() {
            let self = this;

            $.ajax({
                type: "POST",
                url: '/api/bookmarks',
                data: JSON.stringify(this.newBookmark),
                contentType: "application/json",
                success: function () {
                    self.newBookmark = {};
                    self.loadBookmarks();
                }
            });
        },

        deleteBookmark(id) {
            let self = this;
            $.ajax({
                type: "DELETE",
                url: 'api/bookmarks/' + id,
                success: function () {
                    self.loadBookmarks();
                }
            });
        }
    },
    computed: {}
});