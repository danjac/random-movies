<template>
    <div class="container">
    <h1>Random movies</h1>
    <form class="form form-horizontal">
        <div class="form-group">
            <input class="form-control" type="text" placeholder="Find another movie" />
        </div>
        <button class="btn btn-primary form-control" type="submit">
            <span class="glyphicon glyphicon-plus"></span> Add
        </button>
    </form>
    <!-- main view -->
    <router-view
        class="view"
        keep-alive
        transition
        transition-mode="out-in">
    </router-view>
    </div>
</template>
<script>
</script>

<style lang="stylus">
@import '../../node_modules/bootstrap/dist/css/bootstrap.min.css';
</style>


