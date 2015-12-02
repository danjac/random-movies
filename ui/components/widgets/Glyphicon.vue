<template>
    <span class="glyphicon glyphicon-{{icon}}"></span>
</template>

<script>
export default {
    name: "Glyphicon",
    props: {
        icon: {
            type: String,
            required: true
        }
    }
};
</script>
