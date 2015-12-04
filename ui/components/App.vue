<template>
    <div class="container">
    <h1>Random movies</h1>
    <div v-for="alert in alerts" 
         class="alert alert-dismissable alert-{{alert.type}}" 
         role="alert">
        <button type="button" 
                class="close" 
                data-dismiss="alert" 
                aria-label="Close">
            <span aria-hidden="true" @click="removeAlert(alert)">&times;</span>
        </button>
        {{alert.msg}}
    </div>
    <form class="form form-horizontal" @submit="addMovie">
        <div class="form-group">
            <input class="form-control" 
                   type="text" 
                   placeholder="Find another movie" 
                   v-model="title" />
        </div>
        <button class="btn btn-primary form-control" type="submit">
            <glyph icon="plus"></glyph>&nbsp;Add
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
import store from '../store';

export default {
    name: "App",
    data() {
        return {
            title: '',
            alerts: []
        }
    },
    created() {
        store.on('alerts-changed', this.update);
    },
    destroyed() {
        store.removeListener('alerts-changed', this.update);
    },
    methods: {
        update(event) {
            this.alerts = store.getAlerts();
        },
        removeAlert(alert) {
            store.deleteAlert(alert.id);
        },
        addMovie(event) {
            event.preventDefault();
            const title = this.title.trim();
            this.title = "";
            if (!title) {
                return;
            }
            this.$api
            .addMovie(title)
            .then(movie => {
                if (movie) {
                    store.createAlert('New movie added', 'success');
                    this.$route.router.go(`/movie/${movie.imdbID}`);
                }
            })
            .catch(() => {
                store.createAlert("Sorry, couldn't find this movie", 'warning');
            });
        }
    }
};
</script>

<style lang="stylus">
@import '../../node_modules/bootstrap/dist/css/bootstrap.min.css';
</style>


