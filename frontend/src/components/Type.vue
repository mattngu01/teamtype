<script>
export default {
    data() {
        return {
            typedText: ""
        }
    },
    methods: {
    },
    watch: {
        typedText(oldTypedText, newTypedText) {
            // based on assumption that this is called on each char / each change
            if (oldTypedText.length < newTypedText.length) {
                console.log("+" + newTypedText.slice(-1))
            } else {
                console.log("-" + oldTypedText.slice(-1))
            }
        }
    }
}
</script>
<template>
    <div> Hello world </div>
    <input for="test" spellcheck="false" autocapitalize="off" autocomplete="off" v-model="typedText" placeholder="Type here">
</template>